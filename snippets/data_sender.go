package snippets

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type DataSenderConfig struct {
}

type DataSender struct {
	Config DataSenderConfig
}

func (ds *DataSender) Write(ctx *core.Context, t *tuple.Tuple) error {
	return nil
}

func (ds *DataSender) Close(ctx *core.Context) error {
	return nil
}
