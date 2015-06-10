package snippets

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/scouter-snippets/snippets/conf"
	"testing"
	"time"
)

func TestAggregationAndPop(t *testing.T) {
	Convey("Given a integrate mock component", t, func() {
		config := conf.IntegrateConfig{
			FrameInputKeys: []string{"cap1", "cap2"},
		}
		itr := Integrate{
			config:            config,
			trackingInfoQueue: map[string]map[time.Time]trackingInfo{},
		}

		Convey("When aggregate first tracking info", func() {

			now := time.Now()
			ti1 := trackingInfo{
				name:       "from_cap1",
				detectTime: now,
			}
			itr.aggregation(ti1)

			Convey("And when aggregate synchronized tracking info, then will pop 2 tracking info", func() {

				ti2 := trackingInfo{
					name:       "from_cap2",
					detectTime: now,
				}
				itr.aggregation(ti2)
				ok, infos := itr.pop(ti2)
				So(ok, ShouldBeTrue)
				So(len(infos), ShouldEqual, 2)
			})

			Convey("And when not synchronized tracking info, then will not pop", func() {
				ti3 := trackingInfo{
					name:       "from_cap2",
					detectTime: now.Add(1 * time.Second),
				}
				itr.aggregation(ti3)
				ok, infos := itr.pop(ti3)
				So(ok, ShouldBeFalse)
				So(len(infos), ShouldEqual, 0)

				Convey("And when aggregate synchronized next tracking info, then will pop", func() {
					ti4 := trackingInfo{
						name:       "from_cap1",
						detectTime: now.Add(1 * time.Second),
					}
					itr.aggregation(ti4)
					ok, infos := itr.pop(ti4)
					So(ok, ShouldBeTrue)
					So(len(infos), ShouldEqual, 2)
				})
			})
		})
	})
}

func TestAggregationAndPopWithOneInput(t *testing.T) {
	Convey("Given a integrate mock component with one input", t, func() {

		config := conf.IntegrateConfig{
			FrameInputKeys: []string{"cap1"},
		}
		itr := Integrate{
			config:            config,
			trackingInfoQueue: map[string]map[time.Time]trackingInfo{},
		}

		Convey("When aggregate first tracking info", func() {

			now := time.Now()
			ti1 := trackingInfo{
				name:       "from_cap1",
				detectTime: now,
			}
			itr.aggregation(ti1)

			Convey("Then will pop the tracking info", func() {

				ok, infos := itr.pop(ti1)
				So(ok, ShouldBeTrue)
				So(len(infos), ShouldEqual, 1)
			})

			Convey("And when aggregate next tracking info, then will pop it", func() {
				ti2 := trackingInfo{
					name:       "from_cap1",
					detectTime: now.Add(1 * time.Second),
				}
				itr.aggregation(ti2)
				ok, infos := itr.pop(ti2)
				So(ok, ShouldBeTrue)
				So(len(infos), ShouldEqual, 1)
			})
		})
	})
}
