package snippets

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"time"
)

// Tick is source component in order to synchronize all cameras
// frame interval.
type Tick struct {
	tickInterval int64
	finish       bool
}

// SetUp by tick interval. The interval unit is [ms].
func (t *Tick) SetUp(tickInterval int) {
	t.tickInterval = int64(tickInterval)
	t.finish = false
}

// GenerateStream generate tick data in regular interval.
func (t *Tick) GenerateStream(ctx *core.Context, w core.Writer) error {
	for !t.finish {
		select {
		case now := <-time.After(time.Millisecond * time.Duration(t.tickInterval)):
			t := tuple.Tuple{
				Timestamp:     now,
				ProcTimestamp: now,
				Trace:         make([]tuple.TraceEvent, 0),
			}
			w.Write(ctx, &t)
		}
	}
	return nil
}

// Stop generating stream.
func (t *Tick) Stop(ctx *core.Context) error {
	t.finish = true
	return nil
}

// ForcedStop this component
func (t *Tick) ForcedStop() {
	t.finish = true
}
