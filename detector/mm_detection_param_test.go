package detector

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestNewMMDetectionParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a parameter", t, func() {
		params := data.Map{}
		Convey("When the parameter has valid config json file path", func() {
			params["file"] = data.String("mm_detector_param_test.json")
			Convey("Then the state set with detector", func() {
				state, err := createMMDetectionParamState(ctx, params)
				So(err, ShouldBeNil)
				ds, ok := state.(*MMDetectionParamState)
				So(ok, ShouldBeTrue)
				So(ds.d, ShouldNotBeNil)
				ds.d.Delete()
			})
		})
		Convey("When the parameter has invalid param", func() {
			params["filee"] = data.String("mm_detector_param_test.json")
			Convey("Then an error should be occur", func() {
				_, err := createMMDetectionParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has null file path", func() {
			params["file"] = data.Null{}
			Convey("Then an error should be occur", func() {
				_, err := createMMDetectionParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has invalid file path", func() {
			params["file"] = data.String("not_exist.json")
			Convey("Then an error should be occur", func() {
				_, err := createMMDetectionParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
