package plugin

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// PluginSourceCreator is an interface to get bql.SourceCreator.
type PluginSourceCreator interface {
	// CreateSource returns user plug-in source type. Returns error when
	// parameter is invalid.
	CreateSource(ctx *core.Context, with data.Map) (core.Source, error)
	// TypeName return name of registration.
	// Example:
	//  a type name is "capture", then
	//    CREATE SOURCE [source name] TYPE capture WITH ...
	TypeName() string
}
