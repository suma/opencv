package writer

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestVideoWriterCreatorCreatesSink(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}
	removeTestAVIFile()
	Convey("Given a VideoWriter sink creator", t, func() {
		vc := VideoWiterCreator{}
		params := data.Map{}
		Convey("When parameters are empty", func() {
			Convey("Then sink should not be created", func() {
				sink, err := vc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		Convey("When parameter not have file name", func() {
			params["fps"] = data.Float(5)
			params["width"] = data.Int(1920)
			params["height"] = data.Int(1480)
			Convey("Then sink should not be created", func() {
				sink, err := vc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		SkipConvey("When parameter have invalid file name", func() {
			// file name will cast string use `data.ToString`, so any type can
			// cast, and this test is skipped
			params["file_name"] = data.Null{}
			params["fps"] = data.Float(5)
			params["width"] = data.Int(1920)
			params["height"] = data.Int(1480)
			Convey("Then sink should not be created", func() {
				sink, err := vc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		Convey("When parameter have invalid vlaues", func() {
			params["file_name"] = data.String("dummy")
			testMap := data.Map{
				"fps":    data.Null{},
				"width":  data.String("a"),
				"height": data.String("$"),
			}
			for k, v := range testMap {
				msg := fmt.Sprintf("%v error", k)
				Convey("Then sink should not be created because of "+msg, func() {
					params[k] = v
					s, err := vc.CreateSink(ctx, ioParams, params)
					So(err, ShouldNotBeNil)
					So(s, ShouldBeNil)
				})
			}
		})

		Convey("When parameters have only file name", func() {
			params["file_name"] = data.String("dummy")
			Convey("Then sink should be created and other parameters are set default",
				func() {
					sink, err := vc.CreateSink(ctx, ioParams, params)
					defer removeTestAVIFile()
					So(err, ShouldBeNil)
					So(sink, ShouldNotBeNil)
					defer sink.Close(ctx)
					vs, ok := sink.(*videoWriterSink)
					So(ok, ShouldBeTrue)
					So(vs.vw, ShouldNotBeNil)
					_, err = os.Stat("dummy.avi")
					So(os.IsNotExist(err), ShouldBeFalse)
					Convey("And when create another sink with the same file name",
						func() {
							sink2, err := vc.CreateSink(ctx, ioParams, params)
							Convey("Then should not occur an error", func() {
								So(err, ShouldBeNil)
								So(sink2, ShouldNotBeNil)
								defer sink.Close(ctx)
								vs2, ok := sink2.(*videoWriterSink)
								So(ok, ShouldBeTrue)
								So(vs2.vw, ShouldNotBeNil)
							})
						})
				})
		})
	})
}

func removeTestAVIFile() {
	_, err := os.Stat("dummy.avi")
	if !os.IsNotExist(err) {
		os.Remove("dummy.avi")
	}
}
