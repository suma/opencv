package integrator

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core"
)

type InstancesVisualizerFuncCreator struct{}

func (c *InstancesVisualizerFuncCreator) CreateFunction() interface{} {
	return drawWithInstanceStates
}

func (c *InstancesVisualizerFuncCreator) TypeName() string {
	return "draw_with_instance_states"
}

func drawWithInstanceStates(ctx *core.Context, visualizerParam string) (
	[]byte, error) {

	s, err := lookupInstanceVisualizerParamState(ctx, visualizerParam)
	if err != nil {
		return []byte{}, err
	}

	img := s.v.Draw()
	defer img.Delete()

	return img.Serialize(), nil
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
