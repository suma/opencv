package detector

import (
	"encoding/json"
	"io/ioutil"
	"pfi/ComputerVision/scouter-core-conf"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/scouter/utils"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// FrameProcessorParamState is a shared state of frame processor parameter for
// scouter-core.
type FrameProcessorParamState struct {
	fp bridge.FrameProcessor
}

func createFrameProcessorParamState(ctx *core.Context, params data.Map) (core.SharedState,
	error) {
	fpConfig := "{}"
	if p, err := params.Get(utils.FilePath); err == nil {
		path, err := data.AsString(p)
		if err != nil {
			return nil, err
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		fpConfig = string(b)
	} else {
		fpConf := scconf.FrameProcessor{}
		if cp, err := params.Get(utils.CameraParameterFilePath); err == nil {
			cpPath, err := data.AsString(cp)
			if err != nil {
				return nil, err
			}
			cpFile, err := ioutil.ReadFile(cpPath)
			if err != nil {
				return nil, err
			}
			cameraParam := &scconf.CameraParameter{}
			if err := json.Unmarshal(cpFile, &cameraParam); err != nil {
				return nil, err
			}
			fpConf.CameraParameter = cameraParam
		}

		if roip, err := params.Get(utils.ROIParameterFilePath); err == nil {
			roiPath, err := data.AsString(roip)
			if err != nil {
				return nil, err
			}
			roiFile, err := ioutil.ReadFile(roiPath)
			if err != nil {
				return nil, err
			}
			roiParam := &scconf.ROI{}
			if err := json.Unmarshal(roiFile, &roiParam); err != nil {
				return nil, err
			}
			fpConf.ROI = roiParam
		}

		b, err := json.Marshal(fpConf)
		if err != nil {
			return nil, err
		}
		fpConfig = string(b)
	}

	s := &FrameProcessorParamState{}
	s.fp = bridge.NewFrameProcessor(fpConfig)

	return s, nil
}

// CreateNewState creates a state of frame processor parameters. The parameter
// is collected on JSON file, see `scouter::FrameProcessor::Config`, which is
// composition of camera parameters and RIO information. Returns an error when
// cannot read the files.
// This state is updatable, and camera parameter could be update.
//
// Usage of WITH parameters:
//  "file":                  The frame processor file path, include camera
//                           parameter and ROI parameter.
//  "camera_parameter_file": The camera parameter file path.
//  "roi_parameter_file":    The ROI parameter file path.
//
// * if "WITH" parameter includes "file" and other key, the state only load
//   "file" key, others are ignored.
// * these parameter are optional, and the state could be initialized with no
//   "WITH" parameter.
func (s *FrameProcessorParamState) CreateNewState() func(*core.Context, data.Map) (
	core.SharedState, error) {
	return createFrameProcessorParamState
}

// TypeName returns type name.
func (s *FrameProcessorParamState) TypeName() string {
	return "scouter_frame_processor_param"
}

// Terminate the components.
func (s *FrameProcessorParamState) Terminate(ctx *core.Context) error {
	s.fp.Delete()
	return nil
}

// Update the state to reload the JSON file without lock.
//
// Usage of WITH parameters:
//  camera_parameter_file: The file path. Returns an error when cannot read the
//                         file.
func (s *FrameProcessorParamState) Update(params data.Map) error {
	p, err := params.Get(utils.CameraParameterFilePath)
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
