package mjpegserv

import (
	"fmt"
	"io/ioutil"
	"os"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
)

type DebugJPEGSink struct {
	outputDir   string
	jpegQuality int
	detectCount detectCounter
}

type detectCounter struct {
	sync.RWMutex
	count map[string]lockCount
}

type lockCount struct {
	sync.Mutex
	count int
}

func (c *detectCounter) get(k string) (lockCount, bool) {
	c.RLock()
	defer c.RUnlock()
	prev, ok := c.count[k]
	return prev, ok
}

func (c *detectCounter) put(k string, v lockCount) {
	c.Lock()
	defer c.Unlock()
	c.count[k] = v
}

func (s *DebugJPEGSink) Write(ctx *core.Context, t *core.Tuple) error {
	count, err := t.Data.Get("region_count")
	if err != nil {
		return err
	}
	countInt, err := data.ToInt(count)
	if err != nil {
		return err
	}

	name, err := t.Data.Get("name")
	if err != nil {
		return err
	}
	nameStr, err := data.ToString(name)
	if err != nil {
		return err
	}

	if prevCount, ok := s.detectCount.get(nameStr); ok {
		if prevCount.count > int(countInt) {
			ctx.Log().Debug("JPEG has already created")
			return nil
		}
	}
	lc := lockCount{count: int(countInt)}
	s.detectCount.put(nameStr, lc)

	img, err := t.Data.Get("img")
	if err != nil {
		return err
	}
	imgByte, err := data.AsBlob(img)
	if err != nil {
		return err
	}
	imgp := bridge.DeserializeMatVec3b(imgByte)
	defer imgp.Delete()

	fileName := fmt.Sprintf("%v/%v.jpg", s.outputDir, nameStr)
	lc.Lock()
	defer lc.Unlock()
	ioutil.WriteFile(fileName, imgp.ToJpegData(s.jpegQuality), os.ModePerm)
	return nil
}

func (s *DebugJPEGSink) Close(ctx *core.Context) error {
	return nil
}

func (s *DebugJPEGSink) CreateSink(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Sink, error) {
	output, err := params.Get("output")
	if err != nil {
		output = data.String(".")
	}
	outputDir, err := data.AsString(output)
	if err != nil {
		return nil, err
	}

	quality, err := params.Get("quality")
	if err != nil {
		quality = data.Int(50)
	}
	q, err := data.AsInt(quality)
	if err != nil {
		return nil, err
	}

	s.outputDir = outputDir
	s.jpegQuality = int(q)
	s.detectCount = detectCounter{count: map[string]lockCount{}}
	return s, nil
}
