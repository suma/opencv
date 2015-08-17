package recog

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// ImageTaggerCaffeParamState is a shared state used by recognizer.
type ImageTaggerCaffeParamState struct {
	tagger bridge.ImageTaggerCaffe
}

var filePath = data.MustCompilePath("file")

func createImageTaggerCaffeParamState(ctx *core.Context, params data.Map) (
	core.SharedState, error) {
	p, err := params.Get(filePath)
	if err != nil {
		return nil, err
	}
	path, err := data.AsString(p)
	if err != nil {
		return nil, err
	}

	// read file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	taggerConfig := string(b)
	s := &ImageTaggerCaffeParamState{}
	s.tagger = bridge.NewImageTaggerCaffe(taggerConfig)

	return s, nil
}

// CreateNewState creates a state of image tagger by Caffe parameters. The
// parameters is collected on JSON file, see `scouter::ImageTaggerCaffe::Config`,
// which is include caffe model.
//
// Usage of WITH parameter:
//  "file": image tagger by Caffe JSON file path.
func (s *ImageTaggerCaffeParamState) CreateNewState() func(*core.Context, data.Map) (
	core.SharedState, error) {
	return createImageTaggerCaffeParamState
}

// TypeName returns type name.
func (s *ImageTaggerCaffeParamState) TypeName() string {
	return "scouter_image_tagger_caffe"
}

// Terminate the components.
func (s *ImageTaggerCaffeParamState) Terminate(ctx *core.Context) error {
	s.tagger.Delete()
	return nil
}
