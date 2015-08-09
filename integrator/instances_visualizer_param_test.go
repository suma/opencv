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
	ctx.SharedStates.Add("imparam", &InstanceManagerParamState{})
	Convey("Given a parameter", t, func() {
		params := data.Map{}
		Convey("When the parameter has valid configs", func() {
			params["camera_ids"] = data.Int(0)
			params["camera_params"] = data.String("camera_param_test.json")
			params["instance_manager_param"] = data.String("imparam")
			Convey("Then the state should be initialized", func() {
				state, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldBeNil)
				defer state.Terminate(ctx)
				ivs, ok := state.(*InstancesVisualizerParamState)
				So(ok, ShouldBeTrue)
				So(ivs.v, ShouldNotBeNil)
			})
		})

		Convey("When the parameter does not have camera ids", func() {
			params["camera_params"] = data.String("camera_param_test.json")
			params["instance_manager_param"] = data.String("imparam")
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has invalid camera id", func() {
			params["camera_ids"] = data.String("0")
			params["camera_params"] = data.String("camera_param_test.json")
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
			params["camera_ids"] = data.Int(0)
			params["camera_params"] = data.Null{}
			params["instance_manager_param"] = data.String("imparam")
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has not existed camera param path", func() {
			params["camera_ids"] = data.Int(0)
			params["camera_params"] = data.String("not_exist.json")
			params["instance_manager_param"] = data.String("imparam")
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When the parameter does not have im param", func() {
			params["camera_ids"] = data.Int(0)
			params["camera_params"] = data.String("camera_param_test.json")
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has not existed im param", func() {
			params["camera_ids"] = data.Int(0)
			params["camera_params"] = data.String("camera_param_test.json")
			params["instance_manager_param"] = data.String("not_exist")
			Convey("Then an error should be occurred", func() {
				_, err := createInstancesVisualizerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
