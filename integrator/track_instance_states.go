package integrator

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// TrackInstanceStatesUDFCreator is a creator of tracking instance states UDF.
type TrackInstanceStatesUDFCreator struct{}

// CreateFunction returns tracking instance states.
//
// Usage:
//  ```
//  scouter_get_current_instance_states([tracker_param], [manager_param],
//                                      [visualizer_param])
//  ```
//  [tracker_param]
//    * type: string
//    * the state of tracker parameter, detail: scouter_tracker_param
//  [manager_param]
//    * type: string
//    * the state of instance manager parameter,
//      detail: scouter_instance_manager_param
//  [visualizer_param]
//    * type: string
//    * the state of instances visualizer parameter, if set empty then the
//      function not create image with states.
//      detail: scouter_instances_visualizer_param
//
// Return:
//  The function returns `data.Map` including current states and image. Current
//  states ate set with "states" key (type `[]data.Blob`). The image is draw
//  with current states and set with "img" key.
func (c *TrackInstanceStatesUDFCreator) CreateFunction() interface{} {
	return getCurrentInstanceStates
}

// TypeName returns type name.
func (c *TrackInstanceStatesUDFCreator) TypeName() string {
	return "scouter_get_current_instance_states"
}

func getCurrentInstanceStates(ctx *core.Context, trackerParam string,
	instanceManagerParam string, instanceVisualizerParam string) (data.Map, error) {

	trackerState, err := lookupTrackerParamState(ctx, trackerParam)
	if err != nil {
		return nil, err
	}

	managerState, err := lookupInstanceManagerParamState(ctx, instanceManagerParam)
	if err != nil {
		return nil, err
	}

	states := managerState.m.TrackAndGetStates(trackerState.t)
	m := data.Map{}
	if len(states) <= 0 {
		ctx.Log().Debug("instance states is empty")
		m["states"] = data.Array{data.Blob([]byte{})}
	} else {
		defer func() {
			for _, s := range states {
				s.Delete()
			}
		}()

		statesByte := make(data.Array, len(states))
		for i, s := range states {
			statesByte[i] = data.Blob(s.Serialize())
		}

		m["states"] = statesByte
	}

	if instanceVisualizerParam != "" {
		visualizerState, err := lookupInstanceVisualizerParamState(ctx,
			instanceVisualizerParam)
		if err != nil {
			return nil, err
		}
		img := visualizerState.v.DrawWithStates()
		defer img.Delete()
		m["img"] = data.Blob(img.Serialize())
	}

	return m, nil
}

func lookupInstanceManagerParamState(ctx *core.Context, instanceManagerParam string) (
	*InstanceManagerParamState, error) {
	st, err := ctx.SharedStates.Get(instanceManagerParam)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*InstanceManagerParamState); ok {
		return s, nil
	}
	return nil, fmt.Errorf(
		"state '%v' cannot be converted to instance_manager_param.state",
		instanceManagerParam)
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
