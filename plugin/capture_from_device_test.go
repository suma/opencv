package plugin

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
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
			creator, err := capture.GetSourceCreator()
			So(err, ShouldBeNil)
			Convey("Then creator should set capture struct members", func() {
				with := map[string]string{
					"device_id": "0",
					"width":     "500",
					"height":    "600",
					"fps":       "25",
					"camera_id": "101",
				}

				_, err = creator(with)
				So(err, ShouldBeNil)
				So(capture.DeviceID, ShouldEqual, 0)
				So(capture.Width, ShouldEqual, 500)
				So(capture.Height, ShouldEqual, 600)
				So(capture.FPS, ShouldEqual, 25)
				So(capture.CameraID, ShouldEqual, 101)
			})

			Convey("Then creator should occur an error", func() {
				with := map[string]string{
					"width":     "500",
					"height":    "600",
					"fps":       "25",
					"camera_id": "101",
				}

				_, err = creator(with)
				So(err, ShouldNotBeNil)
			})

			Convey("Then creator should set default values", func() {
				with := map[string]string{
					"device_id": "0",
				}

				_, err = creator(with)
				So(err, ShouldBeNil)
				So(capture.DeviceID, ShouldEqual, 0)
				So(capture.Width, ShouldEqual, 0)
				So(capture.Height, ShouldEqual, 0)
				So(capture.FPS, ShouldEqual, 0)
				So(capture.CameraID, ShouldEqual, 0)
			})

			Convey("Then creator should occur parse errors", func() {
				with := map[string]string{
					"device_id": "a",
				}
				_, err = creator(with)
				So(err, ShouldNotBeNil)
			})

			Convey("Then creator should occur parse error on option parameters", func() {
				with := map[string]string{
					"device_id": "0",
				}
				testMap := map[string]string{
					"width":     "a",
					"height":    "b",
					"fps":       "@",
					"camera_id": "#",
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
