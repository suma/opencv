package plugin

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// PluginStateCreator is an interface to get core.SharedSate
type PluginStateCreator interface {
	// NewState returns the SharedState.
	NewState(*core.Context, data.Map) (core.SharedState, error)
	// TypeName returns the SharedState' type.
	TypeName() string
}
