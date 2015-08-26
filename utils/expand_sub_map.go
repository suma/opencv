package utils

import (
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

// ExpandSubMapUDSFCreator is a creator of sub map streaming UDSF.
type ExpandSubMapUDSFCreator struct{}

// CreateStreamFunction returns sub map streaming utility.
func (c *ExpandSubMapUDSFCreator) CreateStreamFunction() interface{} {
	return createExpandSubMapUDSF
}

// TypeName returns type name.
func (c *ExpandSubMapUDSFCreator) TypeName() string {
	return "expand_sub_map"
}

func createExpandSubMapUDSF(ctx *core.Context, decl udf.UDSFDeclarer, stream string,
	subField string) (udf.UDSF, error) {

	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "expand_sub_map",
	}); err != nil {
		return nil, err
	}

	return &expandSubMapUDSF{
		subField: data.MustCompilePath(subField),
	}, nil
}

type expandSubMapUDSF struct {
	subField data.Path
}

func (sf *expandSubMapUDSF) Process(ctx *core.Context, t *core.Tuple,
	w core.Writer) error {
	var submap data.Map
	if value, err := t.Data.Get(sf.subField); err != nil {
		return err
	} else if submap, err = data.AsMap(value); err != nil {
		return err
	}

	traces := []core.TraceEvent{}
	if len(t.Trace) > 0 {
		trace := make([]core.TraceEvent, len(t.Trace), (cap(t.Trace)+1)*2)
		copy(trace, t.Trace)
	}
	tu := &core.Tuple{
		Data:          submap,
		Timestamp:     time.Now(),
		ProcTimestamp: t.ProcTimestamp,
		Trace:         traces,
	}
	return w.Write(ctx, tu)
}

func (sf *expandSubMapUDSF) Terminate(ctx *core.Context) error {
	return nil
}
