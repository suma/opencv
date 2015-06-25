package plugin

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
)

// PluginStateCreator is an interface to get core.SharedSate
type PluginStateCreator interface {
	// NewState returns the SharedState.
	NewState(*core.Context, tuple.Map) (core.SharedState, error)
	// TypeName returns the SharedState' type.
	TypeName() string
	// Func is user specific function: UDF.
	Func(*core.Context, tuple.Value) (tuple.Value, error)
}
