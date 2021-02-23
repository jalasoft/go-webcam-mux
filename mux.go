package wmux

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jalasoft/go-webcam"
)

//NewWebcamMux sreates a new http multiplxer that processes requests on web cameras
func NewWebcamMux() *WebcamMux {
	return new(WebcamMux)
}

var router *mux.Router

func init() {
	router = mux.NewRouter()
	router.HandleFunc("/", allWebcamsHandler)

	devRouter := router.PathPrefix("/dev/{device}").Subrouter()
	devRouter.Handle("/cap", webcamHandlerAdapter{webcamCapability})
	devRouter.Handle("/frm", webcamHandlerAdapter{webcamFramesizes})
	devRouter.Handle("/snap", webcamHandlerAdapter{webcamSnapshot})
	devRouter.Handle("/stream", webcamHandlerAdapter{webcamSnapshotStream})
}

// WebcamMux a multiplexer that implements http.Handler interface
type WebcamMux struct{}

func (s *WebcamMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router.ServeHTTP(w, req)
}

//----------------------------------------------------------------------------
//HANDLER ADAPTER THAT OPENS AND CLOSES WEB CAMERA
//----------------------------------------------------------------------------

type webcamHandlerAdapter struct {
	handleFunc func(webcam.Webcam, http.ResponseWriter, *http.Request)
}

func (handler webcamHandlerAdapter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	deviceName := mux.Vars(req)["device"]
	deviceFile := fmt.Sprintf("/dev/%s", deviceName)

	webcam, err := webcam.OpenWebcam(deviceFile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err = webcam.Close(); err != nil {
			log.Printf("Could not close webcam %s: %v\n", deviceFile, err)
		}
	}()

	handler.handleFunc(webcam, w, req)
}
