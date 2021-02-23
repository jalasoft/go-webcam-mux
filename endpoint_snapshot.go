package wmux

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jalasoft/go-webcam"
)

const (
	DEFAULT_PIXEL_FORMAT = "V4L2_PIX_FMT_MJPEG"
)

func webcamSnapshotStream(camera webcam.Webcam, w http.ResponseWriter, req *http.Request) {

	frmSize, err := getDiscreteFrameSize(camera, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	endChannel := make(chan bool)
	snapChannel := make(chan webcam.Snapshot)
	errChannel := make(chan error)

	w.Header().Add("Content-Type", "multipart/x-mixed-replace; boundary=frame")

	go func() {
		<-req.Context().Done()
		log.Println("Client closed connection")
		endChannel <- true
	}()

	go camera.StreamSnapshots(&frmSize, snapChannel, errChannel, endChannel)

main:
	for {

		select {
		case snap, active := <-snapChannel:
			if active {
				if err := streamSnapshot(w, snap); err != nil {
					log.Println(err.Error())
					break main
				}
			} else {
				break main
			}
		case err := <-errChannel:
			log.Println(err.Error())
			break main
		}
	}
}

func streamSnapshot(w http.ResponseWriter, snap webcam.Snapshot) error {

	length, err := fmt.Fprintf(w, "--frame\nContent-Type: image/jpeg\nContent-Length: %d\n\n", len(snap.Data()))

	if err != nil {
		return err
	}

	length, err = w.Write(snap.Data())

	if err != nil {
		return err
	}

	log.Printf("Sent %d bytes\n", length)

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}

func webcamSnapshot(camera webcam.Webcam, w http.ResponseWriter, req *http.Request) {
	framesize, err := getDiscreteFrameSize(camera, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	snap, err := camera.TakeSnapshot(&framesize)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	length, err := w.Write(snap.Data())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("A new image has been taken, %d bytes\n", length)
}

func getDiscreteFrameSize(camera webcam.Webcam, req *http.Request) (webcam.DiscreteFrameSize, error) {

	selector := camera.DiscreteFrameSize()

	params := req.URL.Query()

	values, present := params["pixfmt"]

	if present && len(values) > 0 {
		pixFmtName := values[0]
		selector.PixelFormatName(pixFmtName)
	}

	values, present = params["w"]

	if present && len(values) > 0 {
		width64, err := strconv.ParseUint(values[0], 10, 32)

		if err != nil {
			return webcam.DiscreteFrameSize{}, err
		}

		log.Printf("Width set explicitly: %d\n", width64)
		selector.Width(uint32(width64))
	}

	values, present = params["h"]

	if present && len(values) > 0 {
		height64, err := strconv.ParseUint(values[0], 10, 32)

		if err != nil {
			return webcam.DiscreteFrameSize{}, err
		}

		log.Printf("Height set explicitly: %d\n", height64)
		selector.Height(uint32(height64))
	}

	return selector.Select()
}
