package plugin

import (
	"pfi/scouter-snippets/snippets"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
)

type SourceCreator func(map[string]string) (core.Source, error)

func register() error {
	creator := func(with map[string]string) (core.Source, error) {
		capture := snippets.Capture{}
		err := capture.SetUp(with["config_path"])
		return &capture, err
	}
	if err := bql.RegisterSourceType("scouter_capture", creator); err != nil {
		return err
	}
	return nil
}
