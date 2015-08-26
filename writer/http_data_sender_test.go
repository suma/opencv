package writer

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"strconv"
	"strings"
	"testing"
)

func TestHttpDataSenderCreatorCreate(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}

	Convey("Given a data sender creator", t, func() {
		sc := HTTPDataSenderCreator{}
		params := data.Map{}
		Convey("When parameters are empty", func() {
			Convey("Then sink should not be created", func() {
				sink, err := sc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		Convey("When parameter not have required value", func() {
			params["host"] = data.String("localhost")
			params["path"] = data.String("")
			Convey("Then sink should not be created", func() {
				sink, err := sc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})

		Convey("When parameter has only port", func() {

			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "test server handle")
			})
			ts := httptest.NewServer(dummyHandler)
			defer ts.Close()
			portstr := strings.Split(ts.URL, ":")[2]
			port, _ := strconv.Atoi(portstr)
			params["port"] = data.Int(port)

			Convey("Then sink should be created", func() {
				sink, err := sc.CreateSink(ctx, ioParams, params)
				So(err, ShouldBeNil)
				So(sink, ShouldNotBeNil)
				defer sink.Close(ctx)

				ds, ok := sink.(*httpDataSenderSink)
				So(ok, ShouldBeTrue)
				So(ds.cli, ShouldNotBeNil)

				Convey("And when a tuple write", func() {

					m := data.Map{
						"a": data.String("homhom"),
					}
					t := &core.Tuple{
						Data: m,
					}
					err = sink.Write(ctx, t)

					Convey("Then sink should send data to dummy server", func() {

						So(err, ShouldBeNil)
					})
				})

			})
		})
	})
}
