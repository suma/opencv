package detector

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestNewCameraParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a parameter", t, func() {
		params := data.Map{}
		Convey("When the parameter has valid config json file path", func() {
			params["file"] = data.String("frame_processor_param_test.json")
			Convey("Then the state set with frame processor", func() {
				state, err := createCameraParamState(ctx, params)
				So(err, ShouldBeNil)
				cs, ok := state.(*CameraParamState)
				So(ok, ShouldBeTrue)
				So(cs.fp, ShouldNotBeNil)
				cs.fp.Delete()
			})
		})
		Convey("When the parameter has invalid param", func() {
			params["filee"] = data.String("frame_processor_param_test.json")
			Convey("Then an error should be occur", func() {
				_, err := createCameraParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has null file path", func() {
			params["file"] = data.Null{}
			Convey("Then an error should be occur", func() {
				_, err := createCameraParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has invalid file path", func() {
			params["file"] = data.String("not_exist.json")
			Convey("Then an error should be occur", func() {
				_, err := createCameraParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestUpdateCameraParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given an initialized state", t, func() {
		params := data.Map{
			"file": data.String("frame_processor_param_test.json"),
		}
		state, err := createCameraParamState(ctx, params)
		So(err, ShouldBeNil)
		cs, ok := state.(*CameraParamState)
		So(ok, ShouldBeTrue)
		defer cs.fp.Delete()
		Convey("When the state is updated with valid config json", func() {
			params2 := data.Map{
				"file": data.String("frame_processor_param_test.json"),
			}
			Convey("Then the state should update and occur no error", func() {
				err := cs.Update(params2)
				So(err, ShouldBeNil)
			})
		})
		Convey("When the state is updated with invalid param", func() {
			params2 := data.Map{}
			Convey("Then an error should be occur", func() {
				err := cs.Update(params2)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state is updated with null param", func() {
			params2 := data.Map{
				"file": data.Null{},
			}
			Convey("Then an error should be occur", func() {
				err := cs.Update(params2)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state is updated with invalid file path", func() {
			params2 := data.Map{
				"file": data.String("not_exist.json"),
			}
			Convey("Then an error should be occur", func() {
				err := cs.Update(params2)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
