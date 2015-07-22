package plugin

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// PluginStateCreator is an interface to get core.SharedSate
type PluginStateCreator interface {
	// CreateNewState returns the SharedState creator function
	CreateNewState() func(*core.Context, data.Map) (core.SharedState, error)
	// TypeName returns the SharedState' type.
	TypeName() string
}
