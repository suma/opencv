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

func createACFDetectionParamState(ctx *core.Context, params data.Map) (core.SharedState, error) {
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

	detectConfig := string(b)
	s := &ACFDetectionParamState{}
	s.d = bridge.NewDetector(detectConfig)

	return s, nil
}

func (s *ACFDetectionParamState) CreateNewState() func(*core.Context, data.Map) (core.SharedState, error) {
	return createACFDetectionParamState
}

func (s *ACFDetectionParamState) TypeName() string {
	return "acf_detection_parameter"
}

func (s *ACFDetectionParamState) Terminate(ctx *core.Context) error {
	s.d.Delete()
	return nil
}
