package utils

import (
	"pfi/sensorbee/sensorbee/core"
	"time"
)

// LogElapseTime write elapse time from start to when this method is called.
// Elapse time unit is millisecond.
func LogElapseTime(ctx *core.Context, place string, start time.Time) {
	end := time.Now()
	elapse := float64(end.Sub(start).Nanoseconds()) / 1e6
	ctx.Log().Debugf("[%v] elapse time[ms]: %.3f", place, elapse)
}
