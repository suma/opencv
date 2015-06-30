package detector

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"testing"
)

func TestValidNewACFDetectionParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a new state", t, func() {
		state := ACFDetectionParamState{}
		Convey("When the state get valid config json", func() {
			with := tuple.Map{
				"file": tuple.String("detector_param_test.json"),
			}
			Convey("Given the state set with detector", func() {
				_, err := state.NewState(ctx, with)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestErrorNewACFDetectionParamState(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a new state", t, func() {
		state := ACFDetectionParamState{}
		Convey("When the state get invalid param", func() {
			with := tuple.Map{}
			Convey("Given an error", func() {
				_, err := state.NewState(ctx, with)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state get null file path", func() {
			with := tuple.Map{
				"file": tuple.Null{},
			}
			Convey("Given an error", func() {
				_, err := state.NewState(ctx, with)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the state get invalid file path", func() {
			with := tuple.Map{
				"file": tuple.String("not_exist.json"),
			}
			Convey("Given an error", func() {
				_, err := state.NewState(ctx, with)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
