package snippets

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

type Tick struct {
	tickInterval int64
	finish       bool
}

// SetUp with tickInterval [ms]
func (t *Tick) SetUp(tickInterval int) {
	t.tickInterval = int64(tickInterval)
	t.finish = false
}

func (t *Tick) GenerateStream(ctx *core.Context, w core.Writer) error {
	temp := time.Now().UnixNano() / int64(time.Millisecond)
	for !t.finish {
		now := time.Now()
		current := now.UnixNano() / int64(time.Millisecond)
		if current-temp > t.tickInterval {
			t := tuple.Tuple{
				Timestamp:     now,
				ProcTimestamp: now, // TODO video capture create time
				Trace:         make([]tuple.TraceEvent, 0),
			}
			w.Write(ctx, &t)
			temp = current
		}
	}
	return nil
}

func (t *Tick) Stop(ctx *core.Context) error {
	t.finish = true
	return nil
}

func (t *Tick) Schema() *core.Schema {
	return nil
}
