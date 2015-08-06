package detector

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// CameraParamState is a shared state of camera parameter for scouter-core.
type CameraParamState struct {
	fp bridge.FrameProcessor
}

func createCameraParamState(ctx *core.Context, params data.Map) (core.SharedState,
	error) {
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

	fpConfig := string(b)
	s := &CameraParamState{}
	s.fp = bridge.NewFrameProcessor(fpConfig)

	return s, nil
}

// CreateNewState creates a state of camera parameters. The parameter is
// collected on JSON file, see `scouter::CameraParameter`. This state is
// updatable.
//
// Usage of WITH parameters:
//  file: The file path. Returns an error when cannot read the file.
func (s *CameraParamState) CreateNewState() func(*core.Context, data.Map) (
	core.SharedState, error) {
	return createCameraParamState
}

func (s *CameraParamState) TypeName() string {
	return "scouter_camera_param"
}

func (s *CameraParamState) Terminate(ctx *core.Context) error {
	s.fp.Delete()
	return nil
}

// Update the state to reload the JSON file without lock.
//
// Usage of IWTH parameters:
//  file: The file path. Returns an error when cannot read the file.
func (s *CameraParamState) Update(params data.Map) error {
	p, err := params.Get("file")
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
