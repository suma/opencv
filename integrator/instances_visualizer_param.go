package integrator

import (
	"encoding/json"
	"io/ioutil"
	"pfi/ComputerVision/scouter-core-conf"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

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

	paths, err := params.Get("camera_parameters")
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

func (s *InstancesVisualizerParamState) CreateNewState() func(*core.Context, data.Map) (
	core.SharedState, error) {
	return createInstancesVisualizerParamState
}

func (s *InstancesVisualizerParamState) TypeName() string {
	return "instances_visualizer_parameter"
}

func (s *InstancesVisualizerParamState) Terminate(ctx *core.Context) error {
	s.v.Delete()
	return nil
}
