package detector

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
			params["file"] = data.String("detector_param_test.json")
			Convey("Then the state set with detector", func() {
				state, err := createACFDetectionParamState(ctx, params)
				So(err, ShouldBeNil)
				ds, ok := state.(*ACFDetectionParamState)
				So(ok, ShouldBeTrue)
				So(ds.d, ShouldNotBeNil)
				ds.d.Delete()
			})
		})
		Convey("When the parameter has invalid param", func() {
			params["filee"] = data.String("detector_param_test.json")
			Convey("Then an error should be occur", func() {
				_, err := createACFDetectionParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has null file path", func() {
			params["file"] = data.Null{}
			Convey("Then an error should be occur", func() {
				_, err := createACFDetectionParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has invalid file path", func() {
			params["file"] = data.String("not_exist.json")
			Convey("Then an error should be occur", func() {
				_, err := createACFDetectionParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestNewAcfDetectionParamWithSeparateFile(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a parameter", t, func() {
		params := data.Map{}
		Convey("When the parameter has valid detection config json file path", func() {
			params["detection_file"] = data.String("detector_exclude_camera_param_test.json")
			Convey("Then the state set with detector", func() {
				state, err := createACFDetectionParamState(ctx, params)
				So(err, ShouldBeNil)
				ds, ok := state.(*ACFDetectionParamState)
				So(ok, ShouldBeTrue)
				So(ds.d, ShouldNotBeNil)
				ds.d.Delete()
			})
			Convey("And when the parameter has valid camera parameter config json file path", func() {
				params["camera_parameter_file"] = data.String("camera_param_test.json")
				Convey("Then the state set with detector", func() {
					state, err := createACFDetectionParamState(ctx, params)
					So(err, ShouldBeNil)
					ds, ok := state.(*ACFDetectionParamState)
					So(ok, ShouldBeTrue)
					So(ds.d, ShouldNotBeNil)
					ds.d.Delete()
				})
			})
			Convey("And when the parameter has invalid camera parameter key", func() {
				params["camera_parameter_filee"] = data.String("camera_param_test.json")
				Convey("Then camera param should be ignored and the state set with detector", func() {
					state, err := createACFDetectionParamState(ctx, params)
					So(err, ShouldBeNil)
					ds, ok := state.(*ACFDetectionParamState)
					So(ok, ShouldBeTrue)
					So(ds.d, ShouldNotBeNil)
					ds.d.Delete()
				})
			})
			Convey("And when the parameter has camera parameter path which is null", func() {
				params["camera_parameter_file"] = data.Null{}
				Convey("Then an error should be occur", func() {
					_, err := createACFDetectionParamState(ctx, params)
					So(err, ShouldNotBeNil)
				})
			})
			Convey("And when the parameter has invalid camera parameter path", func() {
				params["camera_parameter_file"] = data.String("not_exist.json")
				Convey("Then an error should be occur", func() {
					_, err := createACFDetectionParamState(ctx, params)
					So(err, ShouldNotBeNil)
				})
			})
		})
		Convey("When the parameter has invalid detection path key", func() {
			params["detection_filee"] = data.String("detector_exclude_camera_param_test.json")
			Convey("Then an error should be occur", func() {
				_, err := createACFDetectionParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has null detection path", func() {
			params["detection_file"] = data.Null{}
			Convey("Then an error should be occur", func() {
				_, err := createACFDetectionParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has invalid detection path", func() {
			params["detection_file"] = data.String("not_exist.json")
			Convey("Then an error should be occur", func() {
				_, err := createACFDetectionParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
