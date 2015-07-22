package integrator

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestNewACFDetectionParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a parameter", t, func() {
		params := data.Map{}
		Convey("When the parameter has valid config json file path", func() {
			params["file"] = data.String("instance_manager_param_test.json")
			Convey("Then the state set with detector", func() {
				state, err := createInstanceManagerParamState(ctx, params)
				So(err, ShouldBeNil)
				is, ok := state.(*InstanceManagerParamState)
				So(ok, ShouldBeTrue)
				So(is.m, ShouldNotBeNil)
				is.m.Delete()
			})
		})
		Convey("When the parameter has invalid param", func() {
			params["filee"] = data.String("instance_manager_param_test.json")
			Convey("Then an error should be occur", func() {
				_, err := createInstanceManagerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has null file path", func() {
			params["file"] = data.Null{}
			Convey("Then an error should be occur", func() {
				_, err := createInstanceManagerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has invalid file path", func() {
			params["file"] = data.String("not_exist.json")
			Convey("Then an error should be occur", func() {
				_, err := createInstanceManagerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
