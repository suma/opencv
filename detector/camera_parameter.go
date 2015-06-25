package detector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"pfi/ComputerVision/scouter-core-conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
)

type CameraParameterState struct {
	filePath string
}

func (s *CameraParameterState) NewState(ctx *core.Context, with tuple.Map) (core.SharedState, error) {
	p, err := with.Get("file")
	if err != nil {
		return nil, err
	}
	path, err := tuple.AsString(p)
	if err != nil {
		return nil, err
	}

	s.filePath = path
	return s, nil
}

func (s *CameraParameterState) TypeName() string {
	return "camera_parameter"
}

func (s *CameraParameterState) Init(ctx *core.Context) error {
	return nil
}

func (s *CameraParameterState) Write(ctx *core.Context, t *tuple.Tuple) error {
	return nil
}

func (s *CameraParameterState) Terminate(ctx *core.Context) error {
	return nil
}

func (s *CameraParameterState) Func(ctx *core.Context, stateName tuple.Value) (tuple.Value, error) {
	s, err := lookupState(ctx, stateName)
	if err != nil {
		return nil, err
	}

	// read file
	file, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	var cameraParameterConf scconf.CameraParameter
	if err := json.Unmarshal(file, &cameraParameterConf); err != nil {
		return nil, err
	}

	return nil, nil
}

func lookupState(ctx *core.Context, stateName tuple.Value) (*CameraParameterState, error) {
	name, err := tuple.AsString(stateName)
	if err != nil {
		return nil, fmt.Errorf("name of the state must be a string: %v", stateName)
	}

	st, err := ctx.GetSharedState(name)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*CameraParameterState); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to camera_parameter.state", name)
}
