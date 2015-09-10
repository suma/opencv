package integrator

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestNewVisualizerParamState(t *testing.T) {
	cc := &core.ContextConfig{}
	ctx := core.NewContext(cc)
	ctx.SharedStates.Add("imparam", "instance_visualizer_param", &InstanceManagerParamState{})
	Convey("Given a parameter", t, func() {
		params := data.Map{}
		Convey("When the parameter has valid configs", func() {
			params["camera_ids"] = data.Array{data.Int(0)}
			params["camera_parameter_files"] = data.Array{data.String("camera_param_test.json")}
			params["instance_manager_param"] = data.String("imparam")
			Convey("Then the state should be initialized", func() {
				state, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldBeNil)
				ivs, ok := state.(*InstancesVisualizerParamState)
				So(ok, ShouldBeTrue)
				So(ivs.v, ShouldNotBeNil)
				Reset(func() {
					state.Terminate(ctx)
				})
			})
		})

		Convey("When the parameter has empty camera parameters", func() {
			params["instance_manager_param"] = data.String("imparam")
			Convey("Then the state should be initialized with empty parameter", func() {
				state, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldBeNil)
				ivs, ok := state.(*InstancesVisualizerParamState)
				So(ok, ShouldBeTrue)
				So(ivs.v, ShouldNotBeNil)
				Reset(func() {
					state.Terminate(ctx)
				})
			})
		})

		Convey("When the parameter does not have camera ids", func() {
			params["camera_params"] = data.Array{data.String("camera_param_test.json")}
			params["instance_manager_param"] = data.Array{data.String("imparam")}
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has invalid camera id", func() {
			params["camera_ids"] = data.Array{data.String("0")}
			params["camera_parameter_files"] = data.Array{data.String("camera_param_test.json")}
			params["instance_manager_param"] = data.String("imparam")
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When the parameter does not have camera param path", func() {
			params["camera_ids"] = data.Int(0)
			params["instance_manager_param"] = data.String("imparam")
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has invalid camera param path", func() {
			params["camera_ids"] = data.Array{data.Int(0)}
			params["camera_parameter_files"] = data.Array{data.Null{}}
			params["instance_manager_param"] = data.String("imparam")
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has not existed camera param path", func() {
			params["camera_ids"] = data.Array{data.Int(0)}
			params["camera_parameter_files"] = data.Array{data.String("not_exist.json")}
			params["instance_manager_param"] = data.String("imparam")
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When the parameter does not have im param", func() {
			params["camera_ids"] = data.Array{data.Int(0)}
			params["camera_parameter_files"] = data.Array{data.String("camera_param_test.json")}
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has not existed im param", func() {
			params["camera_ids"] = data.Array{data.Int(0)}
			params["camera_parameter_files"] = data.Array{data.String("camera_param_test.json")}
			params["instance_manager_param"] = data.String("not_exist")
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestInstancesVisualizerUpdateCameraParam(t *testing.T) {
	cc := &core.ContextConfig{}
	ctx := core.NewContext(cc)
	ctx.SharedStates.Add("imparam", "instance_visualizer_param", &InstanceManagerParamState{})

	Convey("Given an initialized visualizer state", t, func() {
		params := data.Map{
			"camera_ids":             data.Array{data.Int(0)},
			"camera_parameter_files": data.Array{data.String("camera_param_test.json")},
			"instance_manager_param": data.String("imparam"),
		}
		state, err := createInstancesVisualizerParamState(ctx, params)
		So(err, ShouldBeNil)
		vs, ok := state.(*InstancesVisualizerParamState)
		So(ok, ShouldBeTrue)
		Convey("When the state is updated with valid config json", func() {
			params2 := data.Map{
				"camera_id":             data.Int(0),
				"camera_parameter_file": data.String("camera_param_test.json"),
			}
			Convey("Then the state should update and occur no error", func() {
				err := vs.Update(ctx, params2)
				So(err, ShouldBeNil)
			})
		})

		Convey("When the state is updated with invalid param", func() {
			params2 := data.Map{}
			Convey("Then an error should be occur", func() {
				err := vs.Update(ctx, params2)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state is updated with ID is null param", func() {
			params2 := data.Map{
				"camera_id":             data.Null{},
				"camera_parameter_file": data.String("camera_param_test.json"),
			}
			Convey("Then an error should be occur", func() {
				err := vs.Update(ctx, params2)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state is updated with parameter null param", func() {
			params2 := data.Map{
				"camera_id":             data.Int(0),
				"camera_parameter_file": data.Null{},
			}
			Convey("Then an error should be occur", func() {
				err := vs.Update(ctx, params2)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state is updated with not existed file path", func() {
			params2 := data.Map{
				"camera_id":             data.Int(0),
				"camera_parameter_file": data.String("not_exist.json"),
			}
			Convey("Then an error should be occur", func() {
				err := vs.Update(ctx, params2)
				So(err, ShouldNotBeNil)
			})
		})
		Reset(func() {
			state.Terminate(ctx)
		})
	})
}
