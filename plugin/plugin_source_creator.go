package plugin

import (
	"pfi/sensorbee/sensorbee/bql"
)

// PluginSourceCreator is an interface to get bql.SourceCreator.
type PluginSourceCreator interface {
	// GetSourceCreator returns bql.SourceCreator which is set with
	// parameters (see each Source components godoc).
	GetSourceCreator() bql.SourceCreator
	// TypeName return name of registration.
	// Example:
	//  a type name is "capture", then
	//    CREATE SOURCE [source name] TYPE capture WITH ...
	TypeName() string
}
