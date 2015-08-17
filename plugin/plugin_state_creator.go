package plugin

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// StateCreator is an interface to register the user defined state.
type StateCreator interface {
	// CreateNewState returns the SharedState creator function
	CreateNewState() func(*core.Context, data.Map) (core.SharedState, error)
	// TypeName returns the SharedState' type.
	TypeName() string
}
