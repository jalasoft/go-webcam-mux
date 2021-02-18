package wmux

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jalasoft/go-webcam"
)

//------------------------------------------------------------------------------
//JSON RESOURCES
//------------------------------------------------------------------------------

type discreteFramesizeResource struct {
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
}

type stepwiseFramesizeResource struct {
	MinWidth   uint32 `json:"min_width"`
	MaxWidth   uint32 `json:"max_width"`
	StepWidth  uint32 `json:"step_width"`
	MinHeight  uint32 `json:"min_height"`
	MaxHeight  uint32 `json:"max_height"`
	StepHeight uint32 `json:"step_height"`
}

type pixformatFramesizesResource struct {
	PixFormat            string                      `json:"pix_format"`
	PixFormatDescription string                      `json:"pix_format_description"`
	Discrete             []discreteFramesizeResource `json:"discrete"`
	Stepwise             []stepwiseFramesizeResource `json:"stepwise"`
}

//-------------------------------------------------------------------------------
//ENDPOINT METHOD
//-------------------------------------------------------------------------------

func webcamFramesizes(webcam webcam.Webcam, w http.ResponseWriter, req *http.Request) {
	log.Println("JSEM TU")

	pixFormats, err := webcam.QueryFormats()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := frameSizesForPixformats(webcam, pixFormats)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(result)
}

func frameSizesForPixformats(webcam webcam.Webcam, pixFormats []webcam.PixelFormat) ([]pixformatFramesizesResource, error) {
	result := []pixformatFramesizesResource{}

	for _, pixFormat := range pixFormats {
		sizes, err := pixelFormatSizes(webcam, pixFormat)

		if err != nil {
			return nil, err
		}

		result = append(result, sizes)
	}

	return result, nil
}

func pixelFormatSizes(webcam webcam.Webcam, pixFormat webcam.PixelFormat) (pixformatFramesizesResource, error) {
	sizes, err := webcam.QueryFrameSizes(pixFormat)

	if err != nil {
		return pixformatFramesizesResource{}, err
	}

	discrete := []discreteFramesizeResource{}

	for _, size := range sizes.Discrete() {
		discrete = append(discrete, discreteFramesizeResource{Width: size.Width, Height: size.Height})
	}

	stepwise := []stepwiseFramesizeResource{}

	for _, size := range sizes.Stepwise() {
		resource := stepwiseFramesizeResource{
			MinWidth:   size.MinWidth,
			MaxWidth:   size.MaxWidth,
			StepWidth:  size.StepWidth,
			MinHeight:  size.MinHeight,
			MaxHeight:  size.MaxHeight,
			StepHeight: size.StepHeight,
		}

		stepwise = append(stepwise, resource)
	}

	return pixformatFramesizesResource{
		PixFormat:            pixFormat.Name(),
		PixFormatDescription: pixFormat.Description(),
		Discrete:             discrete,
		Stepwise:             stepwise,
	}, nil
}
