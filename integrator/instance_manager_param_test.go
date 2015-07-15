package integrator

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestValidNewACFDetectionParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a new state", t, func() {
		state := InstanceManagerParamState{}
		Convey("When the state get valid config json", func() {
			with := data.Map{
				"file": data.String("instance_manager_param_test.json"),
			}
			Convey("Then the state set with detector", func() {
				_, err := state.NewState(ctx, with)
				So(err, ShouldBeNil)
				So(state.m, ShouldNotBeNil)
				state.m.Delete()
			})
		})
	})
}

func TestErrorNewACFDetectionParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a new state", t, func() {
		state := InstanceManagerParamState{}
		Convey("When the state get invalid param", func() {
			with := data.Map{}
			Convey("Then an error should be occur", func() {
				_, err := state.NewState(ctx, with)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state get null file path", func() {
			with := data.Map{
				"file": data.Null{},
			}
			Convey("Then an error should be occur", func() {
				_, err := state.NewState(ctx, with)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state get invalid file path", func() {
			with := data.Map{
				"file": data.String("not_exist.json"),
			}
			Convey("Then an error should be occur", func() {
				_, err := state.NewState(ctx, with)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
