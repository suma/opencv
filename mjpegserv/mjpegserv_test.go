package mjpegserv

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"testing"
)

func TestMJPEGServCreateDefaultSink(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a MJPEG server sink", t, func() {
		sink := MJPEGServ{}
		Convey("When parameter is not included port", func() {
			params := tuple.Map{}
			Convey("Then sink has default port number", func() {
				_, err := sink.CreateSink(ctx, params)
				So(err, ShouldBeNil)
				So(sink.port, ShouldEqual, 10090)
			})
		})
		Convey("When parameter has valid port", func() {
			params := tuple.Map{
				"port": tuple.Int(8080),
			}
			Convey("Then sink set port number", func() {
				_, err := sink.CreateSink(ctx, params)
				So(err, ShouldBeNil)
				So(sink.port, ShouldEqual, 8080)
			})
		})
	})
}

func TestMJPEGServCreateSinkWithError(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a MJPEG server sink", t, func() {
		sink := MJPEGServ{}
		Convey("When parameters have an invalid port", func() {
			params := tuple.Map{
				"port": tuple.String("8080"),
			}
			Convey("Then returns an error", func() {
				_, err := sink.CreateSink(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
