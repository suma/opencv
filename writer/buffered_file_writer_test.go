package writer

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestBufferedFileWriterCreatorCreatesSink(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}
	removeTestJSONSFile()
	Convey("Given a buffered writer sink creator", t, func() {
		bwc := BufferedFileWriterCreator{}
		params := data.Map{}
		Convey("When parameters are empty", func() {
			Convey("Then sink should not be created", func() {
				sink, err := bwc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		SkipConvey("When parameter does not have file name", func() {
			// file name will cast string use `data.ToString`, so any type can
			// cast, and this test is skipped
			params["file_name"] = data.Null{}
			Convey("Then sink should not be created", func() {
				sink, err := bwc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})

		Convey("When parameters has valid file name", func() {
			params["file_name"] = data.String("_buffertest/dummy.jsons")
			Convey("Then sink should be created", func() {
				sink, err := bwc.CreateSink(ctx, ioParams, params)
				defer removeTestJSONSFile()
				So(err, ShouldBeNil)
				So(sink, ShouldNotBeNil)
				defer sink.Close(ctx)

				bw, ok := sink.(*bufferedFileSink)
				So(ok, ShouldBeTrue)
				So(bw.f, ShouldNotBeNil)
				_, err = os.Stat("_buffertest/dummy.jsons")
				So(os.IsNotExist(err), ShouldBeFalse)
				Convey("And when write a tuple", func() {
					m := data.Map{
						"a": data.String("homhom"),
						"b": data.Int(0),
					}
					tu := &core.Tuple{
						Data: m,
					}
					err = sink.Write(ctx, tu)
					Convey("Then a tuple map should write down to file", func() {
						So(err, ShouldBeNil)
						b, err := ioutil.ReadFile("_buffertest/dummy.jsons")
						So(err, ShouldBeNil)
						actual := string(b)
						expected := `{"a":"homhom","b":0}
`
						So(actual, ShouldEqual, expected)
						Convey("And when write another tuple", func() {
							m2 := data.Map{
								"c": data.Float(0.1),
								"d": data.Null{},
							}
							tu2 := &core.Tuple{
								Data: m2,
							}
							err := sink.Write(ctx, tu2)
							Convey("Then a tuple map should add to write down to the file", func() {
								So(err, ShouldBeNil)
								b2, err := ioutil.ReadFile("_buffertest/dummy.jsons")
								So(err, ShouldBeNil)
								actual2 := string(b2)
								expected2 := expected + `{"c":0.1,"d":null}
`
								So(actual2, ShouldEqual, expected2)
							})
						})
					})
				})
			})
		})
	})

}

func removeTestJSONSFile() {
	_, err := os.Stat("_buffertest/dummy.jsons")
	if !os.IsNotExist(err) {
		os.Remove("_buffertest/dummy.jsons")
		os.Remove("_buffertest")
	}
}
