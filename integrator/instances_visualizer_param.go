package integrator

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

// InstancesVisualizerParamState is a shared state used by Instance Visualizer.
type InstancesVisualizerParamState struct {
	v bridge.InstancesVisualizer
}

func createInstancesVisualizerParamState(ctx *core.Context, params data.Map) (
	core.SharedState, error) {
	cameraIDInts := []int{}
	if ids, err := params.Get(utils.CameraIDsPath); err == nil {
		cameraIDs, err := data.AsArray(ids)
		if err != nil {
			return nil, err
		}
		cameraIDInts = make([]int, len(cameraIDs))
		for i, id := range cameraIDs {
			idInt, err := data.AsInt(id)
			if err != nil {
				return nil, err
			}
			cameraIDInts[i] = int(idInt)
		}
	}

	// read all file path and convert to camera parameter
	cameraParams := []scconf.CameraParameter{}
	if paths, err := params.Get(utils.CameraParameterFilesPath); err == nil {
		pathsStr, err := data.AsArray(paths)
		if err != nil {
			return nil, err
		}
		cameraParams = make([]scconf.CameraParameter, len(pathsStr))
		for i, path := range pathsStr {
			pathStr, err := data.AsString(path)
			if err != nil {
				return nil, err
			}
			fileByte, err := ioutil.ReadFile(pathStr)
			if err != nil {
				return nil, err
			}

			var param scconf.CameraParameter
			err = json.Unmarshal(fileByte, &param)
			if err != nil {
				return nil, err
			}
			cameraParams[i] = param
		}
	}

	if len(cameraIDInts) != len(cameraParams) {
		return nil, fmt.Errorf("camera ID size and camera parameter file size must be same")
	}

	// make instance visualizer parameter manually
	visualizerParam := scconf.Visualizer{
		CameraIDs:        cameraIDInts,
		CameraParameters: cameraParams,
	}

	b, err := json.Marshal(visualizerParam)
	if err != nil {
		return nil, err
	}
	visualizerConf := string(b)

	imParamName, err := params.Get(utils.InstanceManagerParamPath)
	if err != nil {
		return nil, err
	}
	imParamStr, err := data.AsString(imParamName)
	if err != nil {
		return nil, err
	}
	imState, err := lookupInstanceManagerParamState(ctx, imParamStr)
	if err != nil {
		return nil, err
	}

	s := &InstancesVisualizerParamState{
		v: bridge.NewInstancesVisualizer(&(imState.m), visualizerConf),
	}
	return s, nil
}

// CreateNewState creates a state of Instance Visualizer parameters. The state
// could save or load multi places camera parameters.
//
// Usage of WITH parameter:
//  "camera_ids":             camera IDs ([]int)
//  "camera_parameter_files": camera parameters JSON file paths ([]string)
//  "instance_manager_param": a "scouter_instance_manager_param" UDS name
//
// The order of "camera_ids" and "camera_parameter_files" must correspond with
// each others. For example, the `CREATE STATE` query is
//  * camera_ids=[0, 1]
//  * camera_parameter_files=['file1.json', 'file2.json']
// then this state will save 'file1.json' with ID=0 and 'file2.json' with ID=1.
//
// "instance_amanger_param" is need to create Instance Visualizer instance, so
// need to `CREATE STATE` "scouter_instance_manager_param" UDS.
func (s *InstancesVisualizerParamState) CreateNewState() func(*core.Context, data.Map) (
	core.SharedState, error) {
	return createInstancesVisualizerParamState
}

// TypeName returns type name.
func (s *InstancesVisualizerParamState) TypeName() string {
	return "scouter_instances_visualizer_param"
}

// Terminate the components.
func (s *InstancesVisualizerParamState) Terminate(ctx *core.Context) error {
	s.v.Delete()
	return nil
}

// Update the state to reload the JSON file without global lock. User can update
// projection parameter using camera parameter JSON file with camera ID.
func (s *InstancesVisualizerParamState) Update(params data.Map) error {
	var cameraID int
	if id, err := params.Get(utils.CameraIDPath); err != nil {
		return err
	} else if ci, err := data.AsInt(id); err != nil {
		return err
	} else {
		cameraID = int(ci)
	}

	var path string
	if p, err := params.Get(utils.CameraParameterFilePath); err != nil {
		return err
	} else if path, err = data.AsString(p); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	cpConfig := string(b)
	s.v.UpdateCameraParameter(cameraID, cpConfig)

	return nil
}
