package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"pfi/sensorbee/scouter/utils"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// BufferedFileWriterCreator is a creator of buffered writer.
type BufferedFileWriterCreator struct{}

// CreateSink creates a buffered file writer sink.
func (c *BufferedFileWriterCreator) CreateSink(ctx *core.Context,
	ioParams *bql.IOParams, params data.Map) (core.Sink, error) {

	var fileName string
	if fn, err := params.Get(utils.FileNamePath); err != nil {
		return nil, err
	} else if fileName, err = data.ToString(fn); err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("invalid file path: %v", err.Error())
	}
	dirPath := filepath.Dir(absPath)
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, 0755)
	}

	f, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot open the file path: %v", err.Error())
	}

	return &bufferedFileSink{
		f: f,
	}, nil
}

// TypeName returns name.
func (c *BufferedFileWriterCreator) TypeName() string {
	return "buffered_file_writer"
}

type bufferedFileSink struct {
	f *os.File
}

func (s *bufferedFileSink) Write(ctx *core.Context, t *core.Tuple) error {
	str := t.Data.String()
	if _, err := s.f.WriteString(str + "\n"); err != nil {
		return err
	}
	return nil
}

func (s *bufferedFileSink) Close(ctx *core.Context) error {
	s.f.Close()
	return nil
}
