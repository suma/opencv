package snippets

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

// Tick is source component in order to synchronize all cameras
// frame interval.
type Tick struct {
	tickInterval int64
	finish       bool
}

// SetUp by tick interval [micro second]
func (t *Tick) SetUp(tickInterval int) {
	t.tickInterval = int64(tickInterval)
	t.finish = false
}

// GenerateStream generate tick data in regular interval.
func (t *Tick) GenerateStream(ctx *core.Context, w core.Writer) error {
	temp, _ := tuple.ToInt(tuple.Timestamp(time.Now()))
	for !t.finish {
		now := time.Now()
		current, _ := tuple.ToInt(tuple.Timestamp(now))
		if current-temp > t.tickInterval {
			t := tuple.Tuple{
				Timestamp:     now,
				ProcTimestamp: now,
				Trace:         make([]tuple.TraceEvent, 0),
			}
			w.Write(ctx, &t)
			temp = current
		}
	}
	return nil
}

// Stop generating stream.
func (t *Tick) Stop(ctx *core.Context) error {
	t.finish = true
	return nil
}

// Schema returns registered schema.
func (t *Tick) Schema() *core.Schema {
	return nil
}

// ForcedStop this component
func (t *Tick) ForcedStop() {
	t.finish = true
}
