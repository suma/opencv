package plugin

import (
	"pfi/sensorbee/sensorbee/bql"
)

// PluginSourceCreator is an interface to get bq.SourceCreator.
type PluginSourceCreator interface {
	// GetSourceCreator returns bql.SourceCreator which is set with
	// parameters (see each Source components godoc). When fail to initialize
	// the component, returns error.
	GetSourceCreator() (bql.SourceCreator, error)
	// TypeName return name of registration.
	// Example:
	//  a type name is "capture", then
	//    CREATE SOURCE [source name] TYPE capture WITH ...
	TypeName() string
}
