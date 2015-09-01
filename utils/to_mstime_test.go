package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
	"time"
)

func TestToMSTim(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a timestamp", t, func() {
		no := time.Now()
		now := data.Timestamp(no)
		Convey("When call a function", func() {
			ms, err := toMSTime(ctx, now)
			Convey("Then return integer should be ms time", func() {
				So(err, ShouldBeNil)
				expected := time.Duration(no.UnixNano()) / time.Millisecond
				So(ms, ShouldEqual, int(expected))
			})
		})
	})
}
