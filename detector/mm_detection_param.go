package detector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"pfi/ComputerVision/scouter-core-conf"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/scouter/utils"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// MMDetectionParamState is a shared state used by multi-model detector.
type MMDetectionParamState struct {
	d bridge.MMDetector
}

func createMMDetectionParamState(ctx *core.Context, params data.Map) (
	core.SharedState, error) {
	config := ""
	if p, err := params.Get(utils.FilePath); err == nil {
		path, err := data.AsString(p)
		if err != nil {
			return nil, err
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		config = string(b)
	} else {
		dp, err := params.Get(utils.DetectionFilePath)
		if err != nil {
			return nil, fmt.Errorf(
				"state parameter requires configuration parameter file path")
		}
		detectorFilePath, err := data.AsString(dp)
		if err != nil {
			return nil, err
		}
		detectConfigFile, err := ioutil.ReadFile(detectorFilePath)
		if err != nil {
			return nil, err
		}
		var detectConfig scconf.MultiModelDetector
		err = json.Unmarshal(detectConfigFile, &detectConfig)
		if err != nil {
			return nil, err
		}

		if cp, err := params.Get(utils.CameraParameterFilePath); err == nil {
			cameraParamFilePath, err := data.AsString(cp)
			if err != nil {
				return nil, err
			}
			cameraParamFile, err := ioutil.ReadFile(cameraParamFilePath)
			if err != nil {
				return nil, err
			}
			var cameraParamConfig scconf.CameraParameter
			err = json.Unmarshal(cameraParamFile, &cameraParamConfig)
			if err != nil {
				return nil, err
			}
			detectConfig.CameraParameter = &cameraParamConfig
		}
		b, err := json.Marshal(detectConfig)
		if err != nil {
			return nil, err
		}
		config = string(b)
	}

	if config == "" {
		return nil, fmt.Errorf(
			"state parameter requires configuration parameter file path")
	}
	s := &MMDetectionParamState{}
	s.d = bridge.NewMMDetector(config)

	return s, nil
}

// CreateNewState creates a state of multi-model detector parameters. The
// parameter is collected on JSON file, see `scouter::MMDetector::Config`, which
// is composition of detection.model, camera parameters and so on.
//
// Usage of WITH parameter:
//   "file"          : all detection parameters, include "detection_file" and
//                     "camera_parameter_file" (optional)
//   "detection_file": detection configuration parameters
//   "camera_parameter_file"
//                   : camera parameters (optional)
//
// the state permit blow pattern
// * "file" only
// * "detection_file" only
// * "detection_file" and "camera_parameter_file"
// * if the parameter includes "file" and others key, the state load "file" key
//   only.
func (s *MMDetectionParamState) CreateNewState() func(*core.Context, data.Map) (
	core.SharedState, error) {
	return createMMDetectionParamState
}

// TypeName returns type name.
func (s *MMDetectionParamState) TypeName() string {
	return "scouter_mm_detection_param"
}

// Terminate the components.
func (s *MMDetectionParamState) Terminate(ctx *core.Context) error {
	s.d.Delete()
	return nil
}

// Update the state to reload the JSON file without lock.
//
// Usage of WITH parameters:
//  camera_parameter_file: The camera parameter file path. Returns an error when
//                         cannot read the file.
func (s *MMDetectionParamState) Update(ctx *core.Context, params data.Map) error {
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

	cpConfig := string(b)
	s.d.UpdateCameraParameter(cpConfig)

	return nil
}
