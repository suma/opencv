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

func createMMDetectionParamState(ctx *core.Context, params data.Map) (core.SharedState, error) {
	p, err := params.Get("file")
	if err != nil {
		return nil, err
	}
	path, err := data.AsString(p)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	detectConfig := string(b)
	s := &MMDetectionParamState{}
	s.d = bridge.NewMMDetector(detectConfig)

	return s, nil
}

func (s *MMDetectionParamState) CreateNewState() func(*core.Context, data.Map) (core.SharedState, error) {
	return createMMDetectionParamState
}

func (s *MMDetectionParamState) TypeName() string {
	return "multi_model_detection_parameter"
}

func (s *MMDetectionParamState) Terminate(ctx *core.Context) error {
	s.d.Delete()
	return nil
}
