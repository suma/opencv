package mjpegserv

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestMJPEGServWrite(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}
	params := data.Map{"port": data.Int(8081)}
	img, err := ioutil.ReadFile("test_cvmat")
	if err != nil {
		t.Fatalf("cannot load test image file: %v", err)
	}

	sink := MJPEGServ{}
	if _, err := sink.CreateSink(ctx, ioParams, params); err != nil {
		t.Fatalf("cannot create sink: %v", err)
	}
	Convey("Given a MJPEG server sink and start", t, func() {
		Convey("When passes a tuple", func() {
			m := data.Map{
				"name": data.String("dummy_data"),
				"img":  data.Blob(img),
			}
			tu := &core.Tuple{
				Data: m,
			}
			Convey("Then will pass to the server", func() {
				err := sink.Write(ctx, tu)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestMJPEGServWriteWithInvalidTuple(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}
	params := data.Map{"port": data.Int(8082)}
	sink := MJPEGServ{}
	if _, err := sink.CreateSink(ctx, ioParams, params); err != nil {
		t.Fatalf("cannot create sink: %v", err)
	}
	Convey("Given a MJPEG server sink and start", t, func() {
		Convey("When passes no name tuple", func() {
			m := data.Map{
				"img": data.Blob([]byte{}),
			}
			tu := &core.Tuple{
				Data: m,
			}
			Convey("Then returns an error", func() {
				err := sink.Write(ctx, tu)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When passes invalid name tuple", func() {
			m := data.Map{
				"name": data.Null{},
				"img":  data.Blob([]byte{}),
			}
			tu := &core.Tuple{
				Data: m,
			}
			Convey("Then returns an error", func() {
				err := sink.Write(ctx, tu)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When passes no image tuple", func() {
			m := data.Map{
				"name": data.String("name"),
			}
			tu := &core.Tuple{
				Data: m,
			}
			Convey("Then returns an error", func() {
				err := sink.Write(ctx, tu)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When passes invalid image tuple", func() {
			m := data.Map{
				"name": data.String("name"),
				"img":  data.Int(0),
			}
			tu := &core.Tuple{
				Data: m,
			}
			Convey("Then returns an error", func() {
				err := sink.Write(ctx, tu)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestMJPEGServCreateDefaultSink(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}
	Convey("Given a MJPEG server sink", t, func() {
		sink := MJPEGServ{}
		Convey("When parameter is not included port", func() {
			params := data.Map{}
			Convey("Then sink has default port number", func() {
				_, err := sink.CreateSink(ctx, ioParams, params)
				So(err, ShouldBeNil)
				So(sink.port, ShouldEqual, 10090)
			})
		})
		Convey("When parameter has valid port", func() {
			params := data.Map{
				"port": data.Int(8083),
			}
			Convey("Then sink set port number", func() {
				_, err := sink.CreateSink(ctx, ioParams, params)
				So(err, ShouldBeNil)
				So(sink.port, ShouldEqual, 8083)
			})
		})
	})
}

func TestMJPEGServCreateSinkWithError(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}
	Convey("Given a MJPEG server sink", t, func() {
		sink := MJPEGServ{}
		Convey("When parameters have an invalid port", func() {
			params := data.Map{
				"port": data.String("8080"),
			}
			Convey("Then returns an error", func() {
				_, err := sink.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
