package integrator

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestValidNewTrackerParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a new state", t, func() {
		state := TrackerParamState{}
		Convey("When the state get valid config json", func() {
			param := data.Map{
				"file": data.String("tracker_param_test.json"),
			}
			Convey("Then the state set with tracker", func() {
				_, err := state.NewState(ctx, param)
				So(err, ShouldBeNil)
				So(state.t, ShouldNotBeNil)
				state.t.Delete()
			})
		})
	})
}

func TestErrorNewTrackerParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a new state", t, func() {
		state := TrackerParamState{}
		Convey("When the state get invalid param", func() {
			param := data.Map{}
			Convey("Then an error should be error", func() {
				_, err := state.NewState(ctx, param)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state get null file path", func() {
			param := data.Map{
				"file": data.Null{},
			}
			Convey("Then an error should be occur", func() {
				_, err := state.NewState(ctx, param)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("WHen the state get invalid file path", func() {
			param := data.Map{
				"file": data.String("not_exist.json"),
			}
			Convey("Then an error should be occur", func() {
				_, err := state.NewState(ctx, param)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
