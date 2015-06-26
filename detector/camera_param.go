package detector

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
)

type CameraParameterState struct {
	fp bridge.FrameProcessor
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

	// read file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fpConfig := string(b)
	s.fp = bridge.NewFrameProcessor(fpConfig)

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
	s.fp.Delete()
	return nil
}
