package snippets

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

type CaptureConfig struct {
	Uri string
}

type Capture struct {
	config CaptureConfig
}

func (c *Capture) SetUp(config CaptureConfig) error {
	c.config = config
	return nil
}

func (c *Capture) GenerateStream(ctx *core.Context, w core.Writer) error {
	// TOBE get frames
	now := time.Now()
	t := tuple.Tuple{
		Timestamp:     now,
		ProcTimestamp: now,
		Trace:         make([]tuple.TraceEvent, 0),
	}
	w.Write(ctx, &t)
	return nil
}

func (c *Capture) Schema() *core.Schema {
	return nil
}
