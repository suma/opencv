package writer

import (
	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestMJPEGServWrite(t *testing.T) {
	logger := logrus.New()
	cc := &core.ContextConfig{
		Logger: logger,
	}
	ctx := core.NewContext(cc)
	img, err := ioutil.ReadFile("test_cvmat")
	if err != nil {
		t.Fatalf("cannot load test image file: %v", err)
	}

	Convey("Given a MJPEG server sink and start", t, func() {
		mc := MJPEGServCreator{}
		ioParams := &bql.IOParams{}
		params := data.Map{
			"port": data.Int(8099),
		}
		sink, err := mc.CreateSink(ctx, ioParams, params)
		So(err, ShouldBeNil)
		So(sink, ShouldNotBeNil)
		defer sink.Close(ctx)
		Convey("When passes a tuple", func() {
			m := data.Map{
				"name": data.String("dummy_img_manem"),
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
	Convey("Given a MJPEG server sink and start", t, func() {
		mc := MJPEGServCreator{}
		ioParams := &bql.IOParams{}
		params := data.Map{
			"port": data.Int(8098),
		}
		sink, err := mc.CreateSink(ctx, ioParams, params)
		So(err, ShouldBeNil)
		So(sink, ShouldNotBeNil)
		defer sink.Close(ctx)
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

func TestMJPEGServCreatorCreatesSink(t *testing.T) {
	logger := logrus.New()
	cc := &core.ContextConfig{
		Logger: logger,
	}
	ctx := core.NewContext(cc)
	ioParams := &bql.IOParams{}
	Convey("Given a mjpeg server sink creator", t, func() {
		mc := MJPEGServCreator{}
		params := data.Map{}
		Convey("When parameters do not have port", func() {
			Convey("Then sink should have default port number", func() {
				sink, err := mc.CreateSink(ctx, ioParams, params)
				So(err, ShouldBeNil)
				So(sink, ShouldNotBeNil)
				defer sink.Close(ctx)
				m, ok := sink.(*mjpegServ)
				So(ok, ShouldBeTrue)
				So(m.port, ShouldEqual, 10090)
			})
		})
		Convey("When parameters have invalid port", func() {
			params["port"] = data.String("8097")
			Convey("Then sink should not be initialized and occur an error", func() {
				sink, err := mc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		Convey("When parameters have customized port", func() {
			params["port"] = data.Int(8097)
			Convey("Then sink should have the port number", func() {
				sink, err := mc.CreateSink(ctx, ioParams, params)
				So(err, ShouldBeNil)
				So(sink, ShouldNotBeNil)
				defer sink.Close(ctx)
				m, ok := sink.(*mjpegServ)
				So(ok, ShouldBeTrue)
				So(m.port, ShouldEqual, 8097)
			})
		})
	})
}
