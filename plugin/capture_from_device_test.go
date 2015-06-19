package plugin

import (
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
