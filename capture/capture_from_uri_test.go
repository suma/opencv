package capture

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestGenerateStreamURIError(t *testing.T) {
	ctx := &core.Context{}
	sc := CaptureFromURICreator{}
	ioParams := &bql.IOParams{}
	Convey("Given a CaptureFromURI source with invalid URI", t, func() {
		params := data.Map{
			"uri": data.String("error uri"),
		}
		capture, err := sc.CreateSource(ctx, ioParams, params)
		So(err, ShouldBeNil)
		So(capture, ShouldNotBeNil)
		Convey("When generate stream", func() {
			Convey("Then error has occurred", func() {
				err := capture.GenerateStream(ctx, &dummyWriter{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldStartWith, "error")
			})
		})
	})
}

type dummyWriter struct{}

func (w *dummyWriter) Write(ctx *core.Context, t *core.Tuple) error {
	return nil
}

func TestGetURISourceCreator(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}
	Convey("Given a CaptureFromURI creator", t, func() {
		Convey("When create source with full parameters", func() {
			params := data.Map{
				"uri":              data.String("/data/file.avi"),
				"frame_skip":       data.Int(5),
				"camera_id":        data.Int(1),
				"next_frame_error": data.False,
			}
			Convey("Then creator should initialize capture source", func() {
				s, err := createCaptureFromURI(ctx, ioParams, params)
				So(err, ShouldBeNil)
				capture, ok := s.(*captureFromURI)
				So(ok, ShouldBeTrue)
				So(capture.uri, ShouldEqual, "/data/file.avi")
				So(capture.frameSkip, ShouldEqual, 5)
				So(capture.cameraID, ShouldEqual, 1)
				So(capture.endErrFlag, ShouldBeFalse)
			})
		})

		Convey("When create source with empty uri", func() {
			params := data.Map{
				"frame_skip":       data.Int(5),
				"camera_id":        data.Int(1),
				"next_frame_error": data.False,
			}
			Convey("Then creator should occur an error", func() {
				s, err := createCaptureFromURI(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(s, ShouldBeNil)
			})
		})

		Convey("When create source with only uri", func() {
			params := data.Map{
				"uri": data.String("/data/file.avi"),
			}
			Convey("Then capture should set default values", func() {
				s, err := createCaptureFromURI(ctx, ioParams, params)
				So(err, ShouldBeNil)
				capture, ok := s.(*captureFromURI)
				So(ok, ShouldBeTrue)
				So(capture.uri, ShouldEqual, "/data/file.avi")
				So(capture.frameSkip, ShouldEqual, 0)
				So(capture.cameraID, ShouldEqual, 0)
				So(capture.endErrFlag, ShouldBeTrue)
			})
		})

		Convey("When create source with invalid uri", func() {
			params := data.Map{
				"uri": data.Null{},
			}
			Convey("Then create should occur an error", func() {
				s, err := createCaptureFromURI(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(s, ShouldBeNil)
			})
		})

		Convey("When creatcreate source with invalid option parameters", func() {
			params := data.Map{
				"uri": data.String("/data/file.avi"),
			}
			testMap := data.Map{
				"frame_skip":       data.String("@"),
				"camera_id":        data.String("全角"),
				"next_frame_error": data.String("True"),
			}
			for k, v := range testMap {
				msg := fmt.Sprintf("with %v error", k)
				Convey("Then creator should occur a parse error on option parameters"+msg, func() {
					params[k] = v
					s, err := createCaptureFromURI(ctx, ioParams, params)
					So(err, ShouldNotBeNil)
					So(s, ShouldBeNil)
				})
			}
		})
	})
}
