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
	Convey("Given a CaptureFromURI source with invalid URI", t, func() {
		capture := CaptureFromURI{
			URI: "error uri",
		}
		Convey("When generate stream", func() {
			ctx := core.Context{}
			Convey("Then error has occurred", func() {
				err := capture.GenerateStream(&ctx, &dummyWriter{})
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
	Convey("Given a CaptureFromURI source with", t, func() {
		capture := CaptureFromURI{}
		ioParams := bql.IOParams{}
		Convey("When get source creator", func() {
			creator := capture.CreateSource
			ctx := core.Context{}
			Convey("Then creator should set capture struct members", func() {
				params := data.Map{
					"uri":        data.String("/data/file.avi"),
					"frame_skip": data.Int(5),
					"camera_id":  data.Int(1),
				}

				_, err := creator(&ctx, &ioParams, params)
				So(err, ShouldBeNil)
				So(capture.URI, ShouldEqual, "/data/file.avi")
				So(capture.FrameSkip, ShouldEqual, 5)
				So(capture.CameraID, ShouldEqual, 1)
			})

			Convey("Then creator should occur an error", func() {
				params := data.Map{
					"frame_skip": data.Int(5),
					"camera_id":  data.Int(1),
				}

				_, err := creator(&ctx, &ioParams, params)
				So(err, ShouldNotBeNil)
			})

			Convey("Then creator should set default values", func() {
				params := data.Map{
					"uri": data.String("/data/file.avi"),
				}

				_, err := creator(&ctx, &ioParams, params)
				So(err, ShouldBeNil)
				So(capture.URI, ShouldEqual, "/data/file.avi")
				So(capture.FrameSkip, ShouldEqual, 0)
				So(capture.CameraID, ShouldEqual, 0)
			})

			Convey("Then creator should occur parse error on option parameters", func() {
				params := data.Map{
					"uri": data.String("/data/file.avi"),
				}
				testMap := data.Map{
					"frame_skip": data.String("@"),
					"camera_id":  data.String("全角"),
				}
				for k, v := range testMap {
					Convey(fmt.Sprintf("with %v error", k), func() {
						params[k] = v
						_, err := creator(&ctx, &ioParams, params)
						So(err, ShouldNotBeNil)
					})
				}
			})
		})
	})
}