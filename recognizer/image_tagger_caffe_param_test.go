package recog

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestValidNewImageTaggerCaffeParam(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a new state", t, func() {
		state := ImageTaggerCaffeParamState{}
		Convey("When the state get valid config json", func() {
			params := data.Map{
				"file": data.String("image_tagger_caffe_param_test.json"),
			}
			Convey("Then the state set with tagger", func() {
				_, err := state.NewState(ctx, params)
				So(err, ShouldBeNil)
				So(state.tagger, ShouldNotBeNil)
				state.tagger.Delete()
			})
		})
	})
}

func TestErrorNewImageTaggerCaffeParam(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a new state", t, func() {
		state := ImageTaggerCaffeParamState{}
		Convey("When the state get invalid param", func() {
			params := data.Map{}
			Convey("Then an error should be occur", func() {
				_, err := state.NewState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state get null file path", func() {
			params := data.Map{
				"file": data.Null{},
			}
			Convey("Then an error should be occur", func() {
				_, err := state.NewState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state get invalid file path", func() {
			params := data.Map{
				"file": data.String("not_exist.json"),
			}
			Convey("Then an error should be occur", func() {
				_, err := state.NewState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
