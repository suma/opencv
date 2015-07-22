package integrator

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestNewTrackerParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a parameter", t, func() {
		params := data.Map{}
		Convey("When the parameter has valid config json file path", func() {
			params["file"] = data.String("tracker_param_test.json")
			Convey("Then the state set with tracker", func() {
				state, err := createTrackerParamState(ctx, params)
				So(err, ShouldBeNil)
				ts, ok := state.(*TrackerParamState)
				So(ok, ShouldBeTrue)
				So(ts.t, ShouldNotBeNil)
				ts.t.Delete()
			})
		})
		Convey("When the parameter has invalid param", func() {
			params["filee"] = data.String("tracker_param_test.json")
			Convey("Then an error should be error", func() {
				_, err := createTrackerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has null file path", func() {
			params["file"] = data.Null{}
			Convey("Then an error should be occur", func() {
				_, err := createTrackerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("WHen the parameter has invalid file path", func() {
			params["file"] = data.String("not_exist.json")
			Convey("Then an error should be occur", func() {
				_, err := createTrackerParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
