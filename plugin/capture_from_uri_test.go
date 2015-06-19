package plugin

import (
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
