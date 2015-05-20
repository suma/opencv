package demo

import (
	"pfi/scoutor-snippets/snippets"
	"pfi/sensorbee/sensorbee/core"
)

func demo() (core.StaticTopology, error) {
	tb := core.NewDefaultStaticTopologyBuilder()

	cap1_conf := snippets.CaptureConfig{
		URI: "",
	}
	cap1 := snippets.Capture{}
	cap1.SetUp(cap1_conf)
	tb.AddSource("cap1", &cap1)

	ds_conf := snippets.DetectSimpleConfig{}
	ds := snippets.DetectSimple{
		Config: ds_conf,
	}
	tb.AddBox("detect_simple", &ds).Input("cap1")

	rc_conf := snippets.RecognizeCaffeConfig{}
	rc := snippets.RecognizeCaffe{
		Config: rc_conf,
	}
	tb.AddBox("recognize_caffe", &rc).Input("detect_simple")

	itr_conf := snippets.IntegrateConfig{}
	itr := snippets.Integrate{
		Config: itr_conf,
	}
	tb.AddBox("integrate", &itr).Input("recognize_caffe")

	sender_conf := snippets.DataSenderConfig{}
	sender := snippets.DataSender{
		Config: sender_conf,
	}
	tb.AddSink("data_sender", &sender).Input("integrate")

	return tb.Build()
}
