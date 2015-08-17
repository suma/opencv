package capture

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestGenerateStreamDeviceError(t *testing.T) {
	ctx := &core.Context{}
	sc := FromDeviceCreator{}
	ioParams := &bql.IOParams{}
	Convey("Given a CaptureFromDevice source with invalid device ID", t, func() {
		params := data.Map{
			"device_id": data.Int(999999), // invalid device ID
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

func TestGetDeviceSourceCreator(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}
	Convey("Given a CaptureFromDevice creator", t, func() {
		sc := FromDeviceCreator{}
		Convey("When create source with full parameters", func() {
			params := data.Map{
				"device_id": data.Int(0),
				"width":     data.Int(500),
				"height":    data.Int(600),
				"fps":       data.Int(25),
				"camera_id": data.Int(101),
			}
			Convey("Then creator should initialize capture source", func() {
				s, err := sc.CreateSource(ctx, ioParams, params)
				So(err, ShouldBeNil)
				capture, ok := s.(*captureFromDevice)
				So(ok, ShouldBeTrue)
				So(capture.deviceID, ShouldEqual, 0)
				So(capture.width, ShouldEqual, 500)
				So(capture.height, ShouldEqual, 600)
				So(capture.fps, ShouldEqual, 25)
				So(capture.cameraID, ShouldEqual, 101)
			})
		})

		Convey("When create source with empty device ID", func() {
			params := data.Map{
				"width":     data.Int(500),
				"height":    data.Int(600),
				"fps":       data.Int(25),
				"camera_id": data.Int(101),
			}
			Convey("Then creator should occur an error", func() {
				s, err := sc.CreateSource(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(s, ShouldBeNil)
			})
		})

		Convey("When create source with only device ID", func() {
			params := data.Map{
				"device_id": data.Int(0),
			}
			Convey("Then capture should set default values", func() {
				s, err := sc.CreateSource(ctx, ioParams, params)
				So(err, ShouldBeNil)
				capture, ok := s.(*captureFromDevice)
				So(ok, ShouldBeTrue)
				So(capture.deviceID, ShouldEqual, 0)
				So(capture.width, ShouldEqual, 0)
				So(capture.height, ShouldEqual, 0)
				So(capture.fps, ShouldEqual, 0)
				So(capture.cameraID, ShouldEqual, 0)
			})
		})

		Convey("When create source with invalid device ID", func() {
			params := data.Map{
				"device_id": data.String("a"),
			}
			Convey("Then creator should occur parse errors", func() {
				s, err := sc.CreateSource(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(s, ShouldBeNil)
			})
		})

		Convey("When create source with invalid option parameters", func() {
			params := data.Map{
				"device_id": data.Int(0),
			}
			testMap := data.Map{
				"width":     data.String("a"),
				"height":    data.String("b"),
				"fps":       data.String("@"),
				"camera_id": data.String("#"),
			}
			for k, v := range testMap {
				msg := fmt.Sprintf("with %v error", k)
				Convey("Then creator should occur a parse error on option parameters "+msg,
					func() {
						params[k] = v
						s, err := sc.CreateSource(ctx, ioParams, params)
						So(err, ShouldNotBeNil)
						So(s, ShouldBeNil)
					})
			}
		})
	})
}
