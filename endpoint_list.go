package wmux

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jalasoft/go-webcam"
)

type webcamInfoResource struct {
	Name string `json:"name"`
	File string `json:"file"`
}

func allWebcamsHandler(w http.ResponseWriter, req *http.Request) {

	list, err := webcam.FindWebcams()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resources := []webcamInfoResource{}

	for _, info := range list {
		name := fmt.Sprintf("%s", info.Name)
		resources = append(resources, webcamInfoResource{Name: name, File: info.File.Name()})
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}
