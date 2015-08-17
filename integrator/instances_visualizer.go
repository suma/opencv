package integrator

import (
	"fmt"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type InstancesVisualizerFuncCreator struct{}

func (c *InstancesVisualizerFuncCreator) CreateFunction() interface{} {
	return drawWithInstanceStates
}

func (c *InstancesVisualizerFuncCreator) TypeName() string {
	return "draw_with_instance_states"
}

func drawWithInstanceStates(ctx *core.Context, visualizerParam string,
	frames data.Array, states data.Array, trackees data.Array) (
	[]byte, error) {

	s, err := lookupInstanceVisualizerParamState(ctx, visualizerParam)
	if err != nil {
		return []byte{}, err
	}

	matMap, err := convertToMatVecMap(frames)
	defer func() {
		for _, v := range matMap {
			v.Delete()
		}
	}()
	if err != nil {
		return []byte{}, err
	}

	iss, err := convertToStatesToSlice(states)
	defer func() {
		for _, s := range iss {
			s.Delete()
		}
	}()
	if err != nil {
		return []byte{}, err
	}

	trs, err := convertToTrackeeSlice(trackees)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		for _, tr := range trs {
			tr.MVCandidate.Delete()
		}
	}()

	img := s.v.Draw(matMap, iss, trs)
	defer img.Delete()

	return img.Serialize(), nil
}

var (
	cameraIDPath = data.MustCompilePath("camera_id")
	imgPath      = data.MustCompilePath("img")
)

func convertToMatVecMap(frameArray data.Array) (map[int]bridge.MatVec3b, error) {
	matMap := map[int]bridge.MatVec3b{}
	for _, f := range frameArray {
		fMap, err := data.AsMap(f)
		if err != nil {
			return matMap, err
		}

		id, err := fMap.Get(cameraIDPath)
		if err != nil {
			return matMap, err
		}
		cameraID, err := data.AsInt(id)
		if err != nil {
			return matMap, err
		}

		image, err := fMap.Get(imgPath)
		if err != nil {
			return matMap, err
		}
		imageByte, err := data.AsBlob(image)
		if err != nil {
			return matMap, err
		}

		matMap[int(cameraID)] = bridge.DeserializeMatVec3b(imageByte)
	}
	return matMap, nil
}

func convertToStatesToSlice(states data.Array) ([]bridge.InstanceState, error) {
	iss := []bridge.InstanceState{}
	for _, s := range states {
		sByte, err := data.AsBlob(s)
		if err != nil {
			return iss, err
		}
		state := bridge.DeserializeInstanceState(sByte)

		iss = append(iss, state)
	}
	return iss, nil
}

func lookupInstanceVisualizerParamState(ctx *core.Context, visualizerParam string) (
	*InstancesVisualizerParamState, error) {
	st, err := ctx.SharedStates.Get(visualizerParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*InstancesVisualizerParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf(
		"state '%v' cannot be converted to instance_visualizer_param.state",
		visualizerParam)
}
