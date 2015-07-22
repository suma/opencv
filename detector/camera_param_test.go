package detector

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestValidNewCameraParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a new state", t, func() {
		state := CameraParamState{}
		Convey("When the state get valid config json", func() {
			param := data.Map{
				"file": data.String("frame_processor_param_test.json"),
			}
			Convey("Then the state set with frame processor", func() {
				_, err := state.NewState(ctx, param)
				So(err, ShouldBeNil)
				So(state.fp, ShouldNotBeNil)
				state.fp.Delete()
			})
		})
	})
}

func TestErrorNewCameraParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a new state", t, func() {
		state := CameraParamState{}
		Convey("When the state get invalid param", func() {
			param := data.Map{}
			Convey("Then an error should be occur", func() {
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
		Convey("When the state get invalid file path", func() {
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

func TestUpdateCameraParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given an initialized state", t, func() {
		state := CameraParamState{}
		param := data.Map{
			"file": data.String("frame_processor_param_test.json"),
		}
		_, err := state.NewState(ctx, param)
		So(err, ShouldBeNil)
		defer state.fp.Delete()
		Convey("When the state is updated with valid config json", func() {
			param2 := data.Map{
				"file": data.String("frame_processor_param_test.json"),
			}
			Convey("Then the state should update and occur no error", func() {
				err := state.Update(param2)
				So(err, ShouldBeNil)
			})
		})
		Convey("When the state is updated with invalid param", func() {
			param2 := data.Map{}
			Convey("Then an error should be occur", func() {
				err := state.Update(param2)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state is updated with null param", func() {
			param2 := data.Map{
				"file": data.Null{},
			}
			Convey("Then an error should be occur", func() {
				err := state.Update(param2)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state is updated with invalid file path", func() {
			param2 := data.Map{
				"file": data.String("not_exist.json"),
			}
			Convey("Then an error should be occur", func() {
				err := state.Update(param2)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
