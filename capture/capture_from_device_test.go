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
		ioParams := bql.IOParams{}
		Convey("When get source creator", func() {
			creator := capture.CreateSource
			ctx := core.Context{}
			Convey("Then creator should set capture struct members", func() {
				with := data.Map{
					"device_id": data.Int(0),
					"width":     data.Int(500),
					"height":    data.Int(600),
					"fps":       data.Int(25),
					"camera_id": data.Int(101),
				}

				_, err := creator(&ctx, &ioParams, with)
				So(err, ShouldBeNil)
				So(capture.DeviceID, ShouldEqual, 0)
				So(capture.Width, ShouldEqual, 500)
				So(capture.Height, ShouldEqual, 600)
				So(capture.FPS, ShouldEqual, 25)
				So(capture.CameraID, ShouldEqual, 101)
			})

			Convey("Then creator should occur an error", func() {
				with := data.Map{
					"width":     data.Int(500),
					"height":    data.Int(600),
					"fps":       data.Int(25),
					"camera_id": data.Int(101),
				}

				_, err := creator(&ctx, &ioParams, with)
				So(err, ShouldNotBeNil)
			})

			Convey("Then creator should set default values", func() {
				with := data.Map{
					"device_id": data.Int(0),
				}

				_, err := creator(&ctx, &ioParams, with)
				So(err, ShouldBeNil)
				So(capture.DeviceID, ShouldEqual, 0)
				So(capture.Width, ShouldEqual, 0)
				So(capture.Height, ShouldEqual, 0)
				So(capture.FPS, ShouldEqual, 0)
				So(capture.CameraID, ShouldEqual, 0)
			})

			Convey("Then creator should occur parse errors", func() {
				with := data.Map{
					"device_id": data.String("a"),
				}
				_, err := creator(&ctx, &ioParams, with)
				So(err, ShouldNotBeNil)
			})

			Convey("Then creator should occur parse error on option parameters", func() {
				with := data.Map{
					"device_id": data.Int(0),
				}
				testMap := data.Map{
					"width":     data.String("a"),
					"height":    data.String("b"),
					"fps":       data.String("@"),
					"camera_id": data.String("#"),
				}
				for k, v := range testMap {
					Convey(fmt.Sprintf("with %v error", k), func() {
						with[k] = v
						_, err := creator(&ctx, &ioParams, with)
						So(err, ShouldNotBeNil)
					})
				}
			})
		})
	})
}
