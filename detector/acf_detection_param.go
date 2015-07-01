package detector

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type ACFDetectionParamState struct {
	d bridge.Detector
}

func (s *ACFDetectionParamState) NewState(ctx *core.Context, with data.Map) (core.SharedState, error) {
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
	s.d = bridge.NewDetector(detectConfig)

	return s, nil
}

func (s *ACFDetectionParamState) TypeName() string {
	return "acf_detection_parameter"
}

func (s *ACFDetectionParamState) Init(ctx *core.Context) error {
	return nil
}

func (s *ACFDetectionParamState) Write(ctx *core.Context, t *core.Tuple) error {
	return nil
}

func (s *ACFDetectionParamState) Terminate(ctx *core.Context) error {
	s.d.Delete()
	return nil
}
