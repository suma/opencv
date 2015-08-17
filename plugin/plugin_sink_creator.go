package plugin

import (
	"pfi/sensorbee/sensorbee/bql"
)

// SinkCreator is an interface to register the user defined sink.
type SinkCreator interface {
	// SinkCreator is a function and returns user defined sink type. Returns an
	// error when parameter is invalid.
	bql.SinkCreator
	// TypeName returns name of registration.
	// Example:
	//  a type name is "output_jpeg", then
	//    CREATE SINK [sink name] TYPE output_jpeg WITH ...
	TypeName() string
}
