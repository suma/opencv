package plugin

import (
	"pfi/sensorbee/sensorbee/bql"
)

type PluginSourceCreator interface {
	GetSourceCreator() (bql.SourceCreator, error)
	TypeName() string
}
