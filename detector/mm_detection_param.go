package detector

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type MMDetectionParamState struct {
	d bridge.MMDetector
}

func (s *MMDetectionParamState) NewState(ctx *core.Context, with data.Map) (core.SharedState, error) {
	p, err := with.Get("file")
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

	detectConfig := string(b)
	s.d = bridge.NewMMDetector(detectConfig)

	return s, nil
}

func (s *MMDetectionParamState) TypeName() string {
	return "multi_model_detection_parameter"
}

func (s *MMDetectionParamState) Init(ctx *core.Context) error {
	return nil
}

func (s *MMDetectionParamState) Write(ctx *core.Context, t *core.Tuple) error {
	return nil
}

func (s *MMDetectionParamState) Terminate(ctx *core.Context) error {
	s.d.Delete()
	return nil
}
