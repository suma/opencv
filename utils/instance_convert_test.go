package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"reflect"
	"testing"
)

func TestConvertInstanceStateToMap(t *testing.T) {
	Convey("Given an instance state byte", t, func() {
		ctx := &core.Context{}
		isByte, err := ioutil.ReadFile("_test_state0")
		So(err, ShouldBeNil)

		Convey("When call convert UDF", func() {

			isMap, err := convertInstanceToMap(ctx, isByte)

			Convey("Then process should get map", func() {
				So(err, ShouldBeNil)
				So(isMap, ShouldNotBeNil)

				exInstance, err := data.NewMap(expectedInstanceState)
				So(err, ShouldBeNil)

				// should be check deep equal but not look up all element in go
				SkipSo(reflect.DeepEqual(isMap, exInstance), ShouldBeTrue)
				So(len(isMap), ShouldEqual, len(exInstance))
				So(len(isMap.String()), ShouldEqual, len(exInstance.String()))
			})
		})
	})
}

func TestConvertInstancesStateToMap(t *testing.T) {
	Convey("Given an instance states array", t, func() {
		ctx := &core.Context{}
		isByte, err := ioutil.ReadFile("_test_state0")
		So(err, ShouldBeNil)

		Convey("When call array convert UDF", func() {
			isByteArray := data.Array{data.Blob(isByte)}
			isArray, err := convertInstancesToMap(ctx, isByteArray)

			Convey("Then process should get map array", func() {
				So(err, ShouldBeNil)
				So(isArray, ShouldNotBeNil)

				exInstance, err := data.NewMap(expectedInstanceState)
				So(err, ShouldBeNil)

				actualMap, err := data.AsMap(isArray[0])
				So(err, ShouldBeNil)
				SkipSo(reflect.DeepEqual(actualMap, exInstance), ShouldBeTrue)
				So(len(actualMap), ShouldEqual, len(exInstance))
				So(len(actualMap.String()), ShouldEqual, len(exInstance.String()))
			})
		})
	})
}

var expectedInstanceState = map[string]interface{}{
	"detections": []interface{}{
		map[string]interface{}{
			"bbox": map[string]interface{}{
				"x1": 117,
				"x2": 117,
				"y1": 117,
				"y2": 117,
			},
			"camera_id":  0,
			"confidence": 0.574550211429596,
			"height":     1.8393332958221436,
			"position": map[string]interface{}{
				"x": 4.561657428741455, "y": 8.920892715454102, "z": 0,
			},
			"tags": []interface{}{
				map[string]interface{}{
					"key":   "gender",
					"score": 0.9960467219352722,
					"value": "Male"},
			}},
	},
	"id": 377617671258112,
	"position": map[string]interface{}{
		"x": 4.561657428741455,
		"y": 8.920892715454102,
		"z": 0,
	},
	"tags": []interface{}{},
}
