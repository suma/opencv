package main

import (
	"fmt"
	"pfi/scouter-snippets/snippets"
	"pfi/sensorbee/sensorbee/core"
	"time"
)

func buildTopology() (core.StaticTopology, error) {
	tb := core.NewDefaultStaticTopologyBuilder()

	// TODO use relative URI
	confPath := "/Users/tanakadaisuke/Development/Workspaces/Go/src/pfi/scouter-snippets/configfile/"
	cap1 := snippets.Capture{}
	cap1.SetUp(confPath + "capture[0].json")
	tb.AddSource("cap1", &cap1)

	tick := snippets.Tick{}
	tick.SetUp(50)
	tb.AddSource("tick", &tick)

	ds := snippets.DetectSimple{
		ConfigPath: confPath + "detect[0].json",
	}
	tb.AddBox("detect_simple", &ds).
		NamedInput("cap1", "frame").
		NamedInput("tick", "tick")

	rc := snippets.RecognizeCaffe{
		ConfigPath: confPath + "recognize_caffe[0].json",
	}
	tb.AddBox("recognize_caffe", &rc).Input("detect_simple")
	/*
		itr := snippets.Integrate{
			ConfigPath: confPath + "integrate[0].json",
		}
		tb.AddBox("integrate", &itr).Input("recognize_caffe")

		sender_conf := snippets.DataSenderConfig{}
		sender := snippets.DataSender{
			Config: sender_conf,
		}
		tb.AddSink("data_sender", &sender).Input("integrate")
	*/
	return tb.Build()
}

func main() {
	topoloby, err := buildTopology()
	if err != nil {
		fmt.Printf("topology build error: %v", err.Error())
		return
	}
	logManager := core.NewConsolePrintLogger()
	conf := core.Configuration{
		TupleTraceEnabled: 1,
	}
	ctx := core.Context{
		Logger: logManager,
		Config: conf,
	}
	go topoloby.Run(&ctx)
	time.Sleep(500 * time.Millisecond)
	topoloby.Stop(&ctx)
}
