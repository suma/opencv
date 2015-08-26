package writer

import (
	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestJPEGWiterWrite(t *testing.T) {
	logger := logrus.New()
	cc := &core.ContextConfig{
		Logger: logger,
	}
	ctx := core.NewContext(cc)
	img, err := ioutil.ReadFile("test_cvmat")
	if err != nil {
		t.Fatalf("cannot load test image file: %v", err)
	}
	removeTestJPEGFile()

	Convey("Given a JPEG Writer sink", t, func() {
		jwc := JPEGWriterCreator{}
		ioParams := &bql.IOParams{}
		params := data.Map{}
		sink, err := jwc.CreateSink(ctx, ioParams, params)
		So(err, ShouldBeNil)
		So(sink, ShouldNotBeNil)
		defer sink.Close(ctx)
		Convey("When passes a tuple", func() {
			m := data.Map{
				"name": data.String("test_jpeg"),
				"img":  data.Blob(img),
			}
			tu := &core.Tuple{
				Data: m,
			}
			Convey("Then the tuple should pass through the JPEG writer", func() {
				err := sink.Write(ctx, tu)
				defer removeTestJPEGFile()
				So(err, ShouldBeNil)
				_, err = os.Stat("test_jpeg.jpg")
				So(os.IsNotExist(err), ShouldBeFalse)
			})
		})
	})
}

func removeTestJPEGFile() {
	_, err := os.Stat("test_jpeg.jpg")
	if !os.IsNotExist(err) {
		os.Remove("test_jpeg.jpg")
	}
}

func TestJPEGWriterWithInvalidTuple(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a JPEG Writer sink", t, func() {
		jwc := JPEGWriterCreator{}
		ioParams := &bql.IOParams{}
		params := data.Map{}
		sink, err := jwc.CreateSink(ctx, ioParams, params)
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
		SkipConvey("When passes invalid name tuple", func() {
			// Any type can convert to string, so this case not occur an error
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

func TestJPEGWriterCreatorCreatesSink(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}
	Convey("Given a JPEG writer sink creator", t, func() {
		jwc := JPEGWriterCreator{}
		params := data.Map{}
		Convey("When parameters are empty", func() {
			Convey("Then sink should be created with default parameters", func() {
				sink, err := jwc.CreateSink(ctx, ioParams, params)
				So(err, ShouldBeNil)
				So(sink, ShouldNotBeNil)
				defer sink.Close(ctx)
				js, ok := sink.(*jpegWriterSink)
				So(ok, ShouldBeTrue)
				So(js.outputDir, ShouldEqual, ".")
				So(js.jpegQuality, ShouldEqual, 50)
			})
		})
		Convey("When parameters have invalid directory name", func() {
			params["output"] = data.Null{}
			Convey("Then sink should not be created", func() {
				sink, err := jwc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		Convey("When parameters have invalid JPEG quality", func() {
			params["quality"] = data.String("50")
			Convey("Then sink should not be created", func() {
				sink, err := jwc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		Convey("When parameters set output directory", func() {
			params["output"] = data.String("dummy")
			Convey("Then sink should set output directory name", func() {
				sink, err := jwc.CreateSink(ctx, ioParams, params)
				So(err, ShouldBeNil)
				So(sink, ShouldNotBeNil)
				defer sink.Close(ctx)
				defer removeTestDummyDir()
				js, ok := sink.(*jpegWriterSink)
				So(ok, ShouldBeTrue)
				So(js.outputDir, ShouldEqual, "dummy")
				So(js.jpegQuality, ShouldEqual, 50)
			})
		})
		Convey("When parameters set JPEG quality", func() {
			params["quality"] = data.Int(75)
			Convey("Then sink should set JPEG quality", func() {
				sink, err := jwc.CreateSink(ctx, ioParams, params)
				So(err, ShouldBeNil)
				So(sink, ShouldNotBeNil)
				defer sink.Close(ctx)
				js, ok := sink.(*jpegWriterSink)
				So(ok, ShouldBeTrue)
				So(js.outputDir, ShouldEqual, ".")
				So(js.jpegQuality, ShouldEqual, 75)
			})
		})
	})
}

func removeTestDummyDir() {
	_, err := os.Stat("dummy")
	if !os.IsNotExist(err) {
		os.Remove("dummy")
	}
}
