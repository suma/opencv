package recog

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type ImageTaggerCaffeParamState struct {
	tagger bridge.ImageTaggerCaffe
}

func (s *ImageTaggerCaffeParamState) NewState(ctx *core.Context, params data.Map) (core.SharedState, error) {
	p, err := params.Get("file")
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
	s.tagger = bridge.NewImageTaggerCaffe(taggerConfig)

	return s, nil
}

func (s *ImageTaggerCaffeParamState) TypeName() string {
	return "image_tagger_caffe"
}

func (s *ImageTaggerCaffeParamState) Terminate(ctx *core.Context) error {
	s.tagger.Delete()
	return nil
}
