package wmux

import (
	"encoding/json"
	"net/http"

	"github.com/jalasoft/go-webcam"
)

type capabilitiesResource struct {
	Driver       string   `json:"driver"`
	BusInfo      string   `json:"bus_info"`
	Card         string   `json:"card"`
	Version      uint32   `json:"version"`
	Capabilities []string `json:"capabilities"`
}

func webcamCapability(webcam webcam.Webcam, w http.ResponseWriter, req *http.Request) {

	caps, err := webcam.QueryCapabilities()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	resource := capabilitiesResource{}
	resource.BusInfo = caps.BusInfo()
	resource.Driver = caps.Driver()
	resource.Card = caps.Card()
	resource.Version = caps.Version()
	resource.Capabilities = capsToStrings(caps.Capabilities())

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	encoder.Encode(resource)
}

func capsToStrings(caps []webcam.Capability) []string {
	result := []string{}

	for _, cap := range caps {
		result = append(result, cap.Name)
	}
	return result
}
