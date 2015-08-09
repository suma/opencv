package integrator

import (
	"encoding/json"
	"io/ioutil"
	"pfi/ComputerVision/scouter-core-conf"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// InstancesVisualizerParamState is a shared state used by Instance Visualizer.
type InstancesVisualizerParamState struct {
	v bridge.InstancesVisualizer
}

func createInstancesVisualizerParamState(ctx *core.Context, params data.Map) (
	core.SharedState, error) {
	// camera_ids and camera_parameters should be array type
	ids, err := params.Get("camera_ids")
	if err != nil {
		return nil, err
	}
	cameraIDs, err := data.AsInt(ids)
	if err != nil {
		return nil, err
	}
	cameraIDInts := []int{int(cameraIDs)}

	paths, err := params.Get("camera_params")
	if err != nil {
		return nil, err
	}
	pathStr, err := data.AsString(paths)
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
	cameraParams := []scconf.CameraParameter{param}

	visualizerParam := scconf.Visualizer{
		CameraIDs:        cameraIDInts,
		CameraParameters: cameraParams,
	}

	b, err := json.Marshal(visualizerParam)
	if err != nil {
		return nil, err
	}
	visualizerConf := string(b)

	imParamName, err := params.Get("instance_manager_param")
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
//  "camera_params":          camera parameters JSON file paths ([]string)
//  "instance_manager_param": a "scouter_instance_manager_param" UDS name
//
// The order of "camera_ids" and "camera_params" must correspond with each
// others. For example, the `CREATE STATE` query is
//  * camera_ids=[0, 1]
//  * camera_params=['file1.json', 'file2.json']
// then this state will save 'file1.json' with ID=0 and 'file2.json' with ID=1.
//
// "instance_amanger_param" is need to create Instance Visualizer instance.
func (s *InstancesVisualizerParamState) CreateNewState() func(*core.Context, data.Map) (
	core.SharedState, error) {
	return createInstancesVisualizerParamState
}

func (s *InstancesVisualizerParamState) TypeName() string {
	return "scouter_instances_visualizer_param"
}

func (s *InstancesVisualizerParamState) Terminate(ctx *core.Context) error {
	s.v.Delete()
	return nil
}
