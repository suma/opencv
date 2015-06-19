package plugin

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
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

func (w *dummyWriter) Write(ctx *core.Context, t *tuple.Tuple) error {
	return nil
}

func TestGetURISourceCreator(t *testing.T) {
	Convey("Given a CaptureFromURI source with", t, func() {
		capture := CaptureFromURI{}
		Convey("When get source creator", func() {
			creator, err := capture.GetSourceCreator()
			So(err, ShouldBeNil)
			Convey("Then creator should set capture struct members", func() {
				with := map[string]string{
					"uri":        "/data/file.avi",
					"frame_skip": "5",
					"camera_id":  "1",
				}

				_, err = creator(with)
				So(err, ShouldBeNil)
				So(capture.URI, ShouldEqual, "/data/file.avi")
				So(capture.FrameSkip, ShouldEqual, 5)
				So(capture.CameraID, ShouldEqual, 1)
			})

			Convey("Then creator should occur an error", func() {
				with := map[string]string{
					"frame_skip": "5",
					"camera_id":  "1",
				}

				_, err = creator(with)
				So(err, ShouldNotBeNil)
			})

			Convey("Then creator should set default values", func() {
				with := map[string]string{
					"uri": "/data/file.avi",
				}

				_, err = creator(with)
				So(err, ShouldBeNil)
				So(capture.URI, ShouldEqual, "/data/file.avi")
				So(capture.FrameSkip, ShouldEqual, 0)
				So(capture.CameraID, ShouldEqual, 0)
			})

			Convey("Then creator should occur parse error on option parameters", func() {
				with := map[string]string{
					"uri": "/data/file.avi",
				}
				testMap := map[string]string{
					"frame_skip": "@",
					"camera_id":  "!",
				}
				for k, v := range testMap {
					Convey(fmt.Sprintf("with %v error", k), func() {
						with[k] = v
						_, err = creator(with)
						So(err, ShouldNotBeNil)
					})
				}
			})
		})
	})
}
