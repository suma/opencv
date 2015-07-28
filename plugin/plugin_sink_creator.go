package plugin

import (
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// PluginSinkCreator is an interface to get bql.SinkCreator
type PluginSinkCreator interface {
	// CreateSink returns user plug-in sink type. Returns error when
	// parameter is invalid.
	CreateSink(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Sink, error)
	// TypeName returns name of registration.
	// Example:
	//  a type name is "output_jpeg", then
	//    CREATE SINKE [sink name] TYPE output_jpeg WITH ...
	TypeName() string
}
