package plugin

import (
	"pfi/sensorbee/sensorbee/bql"
)

// SourceCreator is an interface to register the user defined source.
type SourceCreator interface {
	// SourceCreator is a function and returns user defined source. Returns an
	// error when parameter is invalid.
	bql.SourceCreator
	// TypeName returns name of registration.
	// Example:
	//  a type name is "capture", then
	//    CREATE SOURCE [source name] TYPE capture WITH ...
	TypeName() string
}
