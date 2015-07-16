package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestAggregateAllItems(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a aggregation function", t, func() {
		Convey("When execute with an item", func() {
			item1 := data.Int(1)
			actual, err := aggregate(ctx, item1)
			Convey("Then the item should be aggregated to data.Array", func() {
				So(err, ShouldBeNil)
				expected := data.Array{item1}
				So(actual, ShouldResemble, expected)
			})
		})
		Convey("When execute with items which type are all same", func() {
			item1 := data.Int(1)
			item2 := data.Int(2)
			item3 := data.Int(3)
			actual, err := aggregate(ctx, item1, item2, item3)
			Convey("Then these items should be aggregated", func() {
				So(err, ShouldBeNil)
				expected := data.Array{item1, item2, item3}
				So(actual, ShouldResemble, expected)
			})
		})
		Convey("When execute with empty item", func() {
			_, err := aggregate(ctx)
			Convey("Then an error has occur", func() {
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When execute with items which include different type with others", func() {
			item1 := data.Int(1)
			item2 := data.Int(2)
			item3 := data.String("test")
			_, err := aggregate(ctx, item1, item2, item3)
			Convey("Then an error has occur", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
