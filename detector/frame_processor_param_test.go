package detector

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestNewFrameProcessorParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a parameter", t, func() {
		params := data.Map{}
		Convey("When the parameter has valid config json file path", func() {
			params["file"] = data.String("frame_processor_param_test.json")
			Convey("Then the state set with frame processor", func() {
				state, err := createFrameProcessorParamState(ctx, params)
				So(err, ShouldBeNil)
				So(state, ShouldNotBeNil)
				defer state.Terminate(ctx)
				cs, ok := state.(*FrameProcessorParamState)
				So(ok, ShouldBeTrue)
				So(cs.fp, ShouldNotBeNil)
			})
		})

		Convey("When the parameter has 'camera param' and 'ROI param' file paths", func() {
			params["camera_parameter_file"] = data.String("camera_param_test.json")
			params["roi_parameter_file"] = data.String("_roi_param_test.json")
			Convey("Then the state should be set frame processor", func() {
				state, err := createFrameProcessorParamState(ctx, params)
				So(err, ShouldBeNil)
				So(state, ShouldNotBeNil)
				defer state.Terminate(ctx)
				cs, ok := state.(*FrameProcessorParamState)
				So(ok, ShouldBeTrue)
				So(cs.fp, ShouldNotBeNil)
			})
		})

		Convey("When the parameter not have 'file' param", func() {
			params["filee"] = data.String("frame_processor_param_test.json")
			Convey("Then the state set with empty parameter", func() {
				state, err := createFrameProcessorParamState(ctx, params)
				So(err, ShouldBeNil)
				So(state, ShouldNotBeNil)
				defer state.Terminate(ctx)
				cs, ok := state.(*FrameProcessorParamState)
				So(ok, ShouldBeTrue)
				So(cs.fp, ShouldNotBeNil)
			})
		})

		Convey("When the parameter has null file path", func() {
			params["file"] = data.Null{}
			Convey("Then an error should be occur", func() {
				_, err := createFrameProcessorParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has invalid file path", func() {
			params["file"] = data.String("not_exist.json")
			Convey("Then an error should be occur", func() {
				_, err := createFrameProcessorParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has null camera parameter path", func() {
			params["camera_parameter_file"] = data.Null{}
			Convey("Then an error should be occur", func() {
				_, err := createFrameProcessorParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the parameter has not string ROI parameter path", func() {
			params["roi_parameter_file"] = data.Int(0)
			Convey("Then an error should be occur", func() {
				_, err := createFrameProcessorParamState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestUpdateFrameProcessorParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given an initialized state", t, func() {
		params := data.Map{
			"file": data.String("frame_processor_param_test.json"),
		}
		state, err := createFrameProcessorParamState(ctx, params)
		So(err, ShouldBeNil)
		cs, ok := state.(*FrameProcessorParamState)
		So(ok, ShouldBeTrue)
		defer cs.fp.Delete()
		Convey("When the state is updated with valid config json", func() {
			params2 := data.Map{
				"camera_parameter_file": data.String("camera_param_test.json"),
			}
			Convey("Then the state should update and occur no error", func() {
				err := cs.Update(ctx, params2)
				So(err, ShouldBeNil)
			})
		})
		Convey("When the state is updated with invalid param", func() {
			params2 := data.Map{}
			Convey("Then an error should be occur", func() {
				err := cs.Update(ctx, params2)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state is updated with null param", func() {
			params2 := data.Map{
				"camera_parameter_file": data.Null{},
			}
			Convey("Then an error should be occur", func() {
				err := cs.Update(ctx, params2)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state is updated with invalid file path", func() {
			params2 := data.Map{
				"camera_parameter_file": data.String("not_exist.json"),
			}
			Convey("Then an error should be occur", func() {
				err := cs.Update(ctx, params2)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
