package snippets

import (
	"fmt"
	"pfi/scoutor-snippets/snippets/bridge"
	"pfi/scoutor-snippets/snippets/conf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"time"
)

type Capture struct {
	config conf.CaptureConfig
	vcap   bridge.VideoCapture
	fp     bridge.FrameProcessor
}

func (c *Capture) SetUp(configFilePath string) error {
	config, err := conf.GetCaptureSnippetConfig(configFilePath)
	if err != nil {
		return err
	}
	c.config = config
	vcap := bridge.NewVideoCapture()
	if ok := vcap.Open(config.URI); !ok {
		return fmt.Errorf("error opening video stream or file : %v", config.URI)
	}
	c.vcap = vcap

	fp := bridge.NewFrameProcessor(config.FrameProcessorConfig)
	c.fp = fp

	return nil
}

func grab(vcap bridge.VideoCapture, buf bridge.MatVec3b, errChan chan error) {
	if !vcap.IsOpened() {
		errChan <- fmt.Errorf("video stream or file closed")
		return
	}
	tmpBuf := bridge.NewMatVec3b()
	if ok := vcap.Read(tmpBuf); !ok {
		errChan <- fmt.Errorf("cannot read a new frame")
		return
	}
	tmpBuf.CopyTo(buf)
}

func (c *Capture) GenerateStream(ctx *core.Context, w core.Writer) error {
	config := c.config
	rootBuf := bridge.NewMatVec3b()
	buf := bridge.NewMatVec3b()
	var rootBufErr error
	if !config.CaptureFromFile {
		errChan := make(chan error)
		go func(rootBuf bridge.MatVec3b) {
			for {
				grab(c.vcap, rootBuf, errChan)
				select {
				case err := <-errChan:
					rootBufErr = err
					break
				}
			}
		}(rootBuf)
	}

	for { // TODO add stop command and using goroutine
		if config.CaptureFromFile {
			if ok := c.vcap.Read(buf); !ok {
				return fmt.Errorf("cannot read a new frame")
			}
			if config.FrameSkip > 0 {
				for i := 0; i < config.FrameSkip; i++ {
					// TODO considering biding cost
					c.vcap.Grab()
				}
			}
		} else {
			if rootBufErr != nil {
				return rootBufErr
			}
			rootBuf.CopyTo(buf)
			if buf.Empty() {
				return nil //continue
			}
		}

		// TODO create frame processor configuration, very difficult...
		// TODO confirm time stamp using, create in C++ is better?
		now := time.Now()
		inow, _ := tuple.ToInt(tuple.Timestamp(now))
		f := c.fp.Apply(buf, inow, config.CameraID)

		var m = tuple.Map{
			"frame": tuple.Blob(f.Serialize()),
		}
		t := tuple.Tuple{
			Data:          m,
			Timestamp:     now,
			ProcTimestamp: now, // TODO video capture create time
			Trace:         make([]tuple.TraceEvent, 0),
		}
		w.Write(ctx, &t)
	}
	return nil
}

func (c *Capture) Stop(ctx *core.Context) error {
	c.config.FrameProcessorConfig.Delete()
	c.fp.Delete()
	c.vcap.Delete()
	return nil
}

func (c *Capture) Schema() *core.Schema {
	return nil
}
