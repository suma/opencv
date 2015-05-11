package demo

import (
	"pfi/scoutor-snippets/snippets"
	"pfi/sensorbee/sensorbee/core"
)

func demo() (core.StaticTopology, error) {
	tb := core.NewDefaultStaticTopologyBuilder()

	cap1_conf := snippets.CaptureConfig{""}
	cap1 := snippets.Capture{}
	cap1.SetUp(cap1_conf)
	tb.AddSource("cap1", &cap1)

	return tb.Build()
}
