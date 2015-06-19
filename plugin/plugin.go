package plugin

import (
	"pfi/sensorbee/sensorbee/bql"
)

func register() error {
	sources := []PluginSourceCreator{
		&CaptureFromURI{},
		&CaptureFromDevice{},
	}
	for _, source := range sources {
		creator, err := source.GetSourceCreator()
		if err != nil {
			return err
		}
		if err = bql.RegisterSourceType(source.TypeName(), creator); err != nil {
			return err
		}
	}
	return nil
}
