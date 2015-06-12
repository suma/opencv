package snippets

import (
	"pfi/scouter-snippets/snippets/bridge"
	"pfi/scouter-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync"
	"time"
)

// MultiModelDetectSimple detects frames.
type MultiModelDetectSimple struct {
	// ConfigPath is the path of external configuration file
	ConfigPath string

	config    conf.MultiModelDetectSimpleConfig
	detector  bridge.MultiModelDetector
	lastFrame *tuple.Tuple
	mu        sync.RWMutex
}

// Init prepares detection information set by external configuration file.
func (d *MultiModelDetectSimple) Init(ctx *core.Context) error {
	detectConfig, err := conf.GetMultiModelDetectSimpleSnippetConfig(d.ConfigPath)
	if err != nil {
		return err
	}
	d.config = detectConfig
	d.detector = bridge.NewMultiModelDetector(detectConfig.DetectorConfig)
	d.lastFrame = nil
	return nil
}

// Process add detection information to frames. Pass down tuples are controlled by
// tick interval.
func (d *MultiModelDetectSimple) Process(ctx *core.Context, t *tuple.Tuple, w core.Writer) error {
	switch t.InputName {
	case "frame":
		d.mu.Lock()
		defer d.mu.Unlock()
		d.lastFrame = t

	case "tick":
		lastFrame := d.getLastFrame()
		if lastFrame == nil {
			return nil
		}
		d.mu.Lock()
		defer d.mu.Unlock()
		d.lastFrame = nil
		err := d.detect(lastFrame, t.Timestamp)
		if err != nil {
			return err
		}

		w.Write(ctx, lastFrame)
	}
	return nil
}

func (d *MultiModelDetectSimple) getLastFrame() *tuple.Tuple {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.lastFrame == nil {
		return nil
	}
	return d.lastFrame.Copy()
}

func (d *MultiModelDetectSimple) detect(t *tuple.Tuple, timestamp time.Time) error {
	f, err := getFrame(t)
	if err != nil {
		return err
	}

	fPointer := bridge.DeserializeFrame(f)
	defer fPointer.Delete()
	s, _ := tuple.ToInt(tuple.Timestamp(time.Now()))

	drPointer := d.detector.Detect(fPointer)

	t.Data["detection_result"] = tuple.Blob(drPointer.Serialize())
	t.Data["detection_time"] = tuple.Timestamp(timestamp)

	if d.config.PlayerFlag {
		e, _ := tuple.ToInt(tuple.Timestamp(time.Now()))
		ms := e - s
		drw := bridge.DetectDrawResult(fPointer, drPointer, ms)
		defer drw.Delete()
		t.Data["detection_draw_result"] = tuple.Blob(drw.ToJpegData(d.config.JpegQuality))
	}
	return nil
}

// Terminate this component.
func (d *MultiModelDetectSimple) Terminate(ctx *core.Context) error {
	d.detector.Delete()
	return nil
}
