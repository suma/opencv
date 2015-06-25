package capture

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"testing"
)

func TestGenerateStreamDeviceError(t *testing.T) {
	Convey("Given a CaptureFromDevice source with invalid device ID", t, func() {
		capture := CaptureFromDevice{
			DeviceID: 999999,
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

func TestGetDeviceSourceCreator(t *testing.T) {
	Convey("Given a CaptureFromDevice source with", t, func() {
		capture := CaptureFromDevice{}
		Convey("When get source creator", func() {
			creator := capture.GetSourceCreator()
			ctx := core.Context{}
			Convey("Then creator should set capture struct members", func() {
				with := tuple.Map{
					"device_id": tuple.Int(0),
					"width":     tuple.Int(500),
					"height":    tuple.Int(600),
					"fps":       tuple.Int(25),
					"camera_id": tuple.Int(101),
				}

				_, err := creator(&ctx, with)
				So(err, ShouldBeNil)
				So(capture.DeviceID, ShouldEqual, 0)
				So(capture.Width, ShouldEqual, 500)
				So(capture.Height, ShouldEqual, 600)
				So(capture.FPS, ShouldEqual, 25)
				So(capture.CameraID, ShouldEqual, 101)
			})

			Convey("Then creator should occur an error", func() {
				with := tuple.Map{
					"width":     tuple.Int(500),
					"height":    tuple.Int(600),
					"fps":       tuple.Int(25),
					"camera_id": tuple.Int(101),
				}

				_, err := creator(&ctx, with)
				So(err, ShouldNotBeNil)
			})

			Convey("Then creator should set default values", func() {
				with := tuple.Map{
					"device_id": tuple.Int(0),
				}

				_, err := creator(&ctx, with)
				So(err, ShouldBeNil)
				So(capture.DeviceID, ShouldEqual, 0)
				So(capture.Width, ShouldEqual, 0)
				So(capture.Height, ShouldEqual, 0)
				So(capture.FPS, ShouldEqual, 0)
				So(capture.CameraID, ShouldEqual, 0)
			})

			Convey("Then creator should occur parse errors", func() {
				with := tuple.Map{
					"device_id": tuple.String("a"),
				}
				_, err := creator(&ctx, with)
				So(err, ShouldNotBeNil)
			})

			Convey("Then creator should occur parse error on option parameters", func() {
				with := tuple.Map{
					"device_id": tuple.Int(0),
				}
				testMap := tuple.Map{
					"width":     tuple.String("a"),
					"height":    tuple.String("b"),
					"fps":       tuple.String("@"),
					"camera_id": tuple.String("#"),
				}
				for k, v := range testMap {
					Convey(fmt.Sprintf("with %v error", k), func() {
						with[k] = v
						_, err := creator(&ctx, with)
						So(err, ShouldNotBeNil)
					})
				}
			})
		})
	})
}
