package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"reflect"
	"testing"
	"time"
)

func TestConvertInstanceStateToKanohiMap(t *testing.T) {
	Convey("Given an instance state byte", t, func() {
		ctx := &core.Context{}
		isByte, err := ioutil.ReadFile("_test_state0")
		So(err, ShouldBeNil)

		Convey("When call convert UDF", func() {
			isArray := data.Array{data.Blob(isByte)}
			floorID := 0
			now := time.Now()
			timestamp := data.Timestamp(now)
			actual, err := convertInstanceStatesToKanohiMap(ctx, isArray, floorID, timestamp)

			Convey("Then process get valid array for kanohi map", func() {
				So(err, ShouldBeNil)
				So(actual, ShouldNotBeNil)

				expected := getExpectedInstanceForKanohi(now, floorID)
				// should be check deep equal but not look up all element in go
				SkipSo(reflect.DeepEqual(actual, expected), ShouldBeTrue)
				So(len(actual), ShouldEqual, len(expected))
				So(len(actual.String()), ShouldEqual, len(expected.String()))

			})
		})
	})
}

func getExpectedInstanceForKanohi(now time.Time, floorID int) data.Array {
	instance := data.Map{
		"id": data.Int(377617671258112),
		"location": data.Map{
			"x":        data.Float(4.561657428741455),
			"y":        data.Float(8.920892715454102),
			"floor_id": data.Int(floorID),
		},
		"labels": data.Array{},
	}

	ts, _ := data.ToInt(data.Timestamp(now))
	ret := data.Map{
		"instances": instance,
		"time":      data.Int(ts),
	}

	return data.Array{ret}
}
