package detector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"pfi/ComputerVision/scouter-core-conf"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type ACFDetectionParamState struct {
	d bridge.Detector
}

// createACFDetectionParamState creates the core.SharedState for
// multi model detector, which has detector instance.
//
// WITH parameter:
//   "file": all detection parameters, include "detection_file" and
//           "camera_parameter_file" (optional)
//   "detection_file": detection configuration parameters
//   "camera_parameter_file": camera parameters (optional)
//
// the state permit blow pattern
// * "file" only
// * "detection_file" only
// * "detection_file" and "camera_parameter_file"
// * if the parameter has "file" and others key, the state "file" key only.
func createACFDetectionParamState(ctx *core.Context, params data.Map) (core.SharedState, error) {
	config := ""
	if p, ok := params["file"]; ok {
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
		dp, ok := params["detection_file"]
		if !ok {
			return nil, fmt.Errorf("state parameter requires configuration parameter file path")
		}
		detectorFilePath, err := data.AsString(dp)
		if err != nil {
			return nil, err
		}
		detectConfigFile, err := ioutil.ReadFile(detectorFilePath)
		if err != nil {
			return nil, err
		}
		var detectConfig scconf.Detector
		err = json.Unmarshal(detectConfigFile, &detectConfig)
		if err != nil {
			return nil, err
		}

		if cp, ok := params["camera_parameter_file"]; ok {
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
		return nil, fmt.Errorf("state parameter requires configuration parameter file path")
	}
	s := &ACFDetectionParamState{}
	s.d = bridge.NewDetector(config)

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

func (s *ACFDetectionParamState) Update(params data.Map) error {
	p, err := params.Get("camera_parameter_file")
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
