package detector

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type CameraParamState struct {
	fp bridge.FrameProcessor
}

func (s *CameraParamState) NewState(ctx *core.Context, param data.Map) (core.SharedState, error) {
	p, err := param.Get("file")
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

	fpConfig := string(b)
	s.fp = bridge.NewFrameProcessor(fpConfig)

	return s, nil
}

func (s *CameraParamState) TypeName() string {
	return "camera_parameter"
}

func (s *CameraParamState) Terminate(ctx *core.Context) error {
	s.fp.Delete()
	return nil
}

func (s *CameraParamState) Update(param data.Map) error {
	p, err := param.Get("file")
	if err != nil {
		return err
	}
	path, err := data.AsString(p)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	fpConfig := string(b)
	s.fp.UpdateConfig(fpConfig)

	return nil
}
