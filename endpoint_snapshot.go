package wmux

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/jalasoft/go-webcam"
)

const (
	DEFAULT_PIXEL_FORMAT = "V4L2_PIX_FMT_MJPEG"
)

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

	pixFormat, err := getPixelFormat(camera, req)

	if err != nil {
		return webcam.DiscreteFrameSize{}, err
	}

	width, height, err := getWidthHeight(camera, pixFormat, req)

	if err != nil {
		return webcam.DiscreteFrameSize{}, nil
	}

	return webcam.DiscreteFrameSize{
		PixelFormat: pixFormat,
		Width:       width,
		Height:      height,
	}, nil
}

func getPixelFormat(webcam webcam.Webcam, req *http.Request) (webcam.PixelFormat, error) {

	formats, err := webcam.QueryFormats()

	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	values, present := params["pixfmt"]

	if !present || len(values) == 0 {
		return getDefaultFormat(formats), nil
	}

	value := values[0]

	return getFormatByName(formats, value), nil

}

func getDefaultFormat(formats []webcam.PixelFormat) webcam.PixelFormat {

	for _, value := range formats {
		if value.Name() == DEFAULT_PIXEL_FORMAT {
			return value
		}
	}

	return formats[0]
}

func getFormatByName(formats []webcam.PixelFormat, pixFormatName string) webcam.PixelFormat {

	for _, value := range formats {
		if value.Name() == pixFormatName {
			return value
		}
	}

	return getDefaultFormat(formats)
}

func getWidthHeight(camera webcam.Webcam, pixFormat webcam.PixelFormat, req *http.Request) (uint32, uint32, error) {

	frmSizes, err := camera.QueryFrameSizes(pixFormat)

	if err != nil {
		return 0, 0, err
	}

	discrete := frmSizes.Discrete()

	width, err := getWidth(discrete, req)

	if err != nil {
		return 0, 0, err
	}

	matchingHeights := []webcam.DiscreteFrameSize{}

	for _, size := range discrete {
		if size.Width == width {
			matchingHeights = append(matchingHeights, size)
		}
	}

	height, err := getHeight(matchingHeights, req)

	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func getWidth(frmSizes []webcam.DiscreteFrameSize, req *http.Request) (uint32, error) {
	params := req.URL.Query()

	values, present := params["w"]

	if !present || len(values) == 0 {
		return frmSizes[0].Width, nil
	}

	strValue := values[0]
	value64, err := strconv.ParseUint(strValue, 10, 32)

	if err != nil {
		return 0, err
	}

	value32 := uint32(value64)
	min_diff := float64(100000)
	best_width := uint32(0)

	for _, size := range frmSizes {
		diff := math.Abs(float64(value32 - size.Width))
		if diff < min_diff {
			min_diff = diff
			best_width = size.Width
		}
	}
	return best_width, nil
}

func getHeight(frmSizes []webcam.DiscreteFrameSize, req *http.Request) (uint32, error) {

	params := req.URL.Query()

	values, present := params["h"]

	if !present || len(values) == 0 {
		return frmSizes[0].Height, nil
	}

	strValue := values[0]
	value64, err := strconv.ParseUint(strValue, 10, 32)
	value32 := uint32(value64)

	if err != nil {
		return 0, err
	}

	min_diff := float64(100000)
	best_width := uint32(0)

	for _, size := range frmSizes {
		diff := math.Abs(float64(value32 - size.Height))
		if diff < min_diff {
			min_diff = diff
			best_width = size.Width
		}
	}
	return best_width, nil
}
