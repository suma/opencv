package utils

import (
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

// ArrayToStreamUDSFCreator is a creator of split array and generate stream UDSF.
type ArrayToStreamUDSFCreator struct{}

// CreateStreamFunction returns split array and stream utility.
func (c *ArrayToStreamUDSFCreator) CreateStreamFunction() interface{} {
	return createArrayToStreamUDSF
}

// TypeName returns type name.
func (c *ArrayToStreamUDSFCreator) TypeName() string {
	return "array_to_stream"
}

func createArrayToStreamUDSF(ctx *core.Context, decl udf.UDSFDeclarer,
	stream string, arrayName string, outName string) (udf.UDSF, error) {

	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "array_to_stream",
	}); err != nil {
		return nil, err
	}

	if arrayName == "" {
		arrayName = "array"
	}
	if outName == "" {
		outName = "value"
	}

	return &arrayToStreamUDSF{
		arrayName: data.MustCompilePath(arrayName),
		outName:   outName,
	}, nil
}

type arrayToStreamUDSF struct {
	arrayName data.Path
	outName   string
}

func (sf *arrayToStreamUDSF) Process(ctx *core.Context, t *core.Tuple,
	w core.Writer) error {

	var array data.Array
	if a, err := t.Data.Get(sf.arrayName); err != nil {
		return err
	} else if array, err = data.AsArray(a); err != nil {
		return err
	}

	traceCopyFlag := len(t.Trace) > 0
	for _, v := range array {
		now := time.Now()
		m := data.Map{
			sf.outName: v,
		}
		traces := []core.TraceEvent{}
		if traceCopyFlag { // reduce copy cost when trace mode is off
			traces := make([]core.TraceEvent, len(t.Trace), (cap(t.Trace)+1)*2)
			copy(traces, t.Trace)
		}
		tu := &core.Tuple{
			Data:          m,
			Timestamp:     now,
			ProcTimestamp: t.ProcTimestamp,
			Trace:         traces,
		}
		w.Write(ctx, tu)
	}
	return nil
}

func (sf *arrayToStreamUDSF) Terminate(ctx *core.Context) error {
	return nil
}
