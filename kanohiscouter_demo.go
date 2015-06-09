package main

import (
	"fmt"
	"pfi/scouter-snippets/snippets"
	"pfi/sensorbee/sensorbee/core"
	"runtime"
	"time"
)

var (
	cameras []*snippets.Capture = []*snippets.Capture{}
	ticker  *snippets.Tick
)

func buildTopology() (core.StaticTopology, error) {
	tb := core.NewDefaultStaticTopologyBuilder()

	confPath := "./configfile/"

	cap1 := snippets.Capture{}
	cap1.SetUp(confPath + "capture[0].json")
	tb.AddSource("cap1", &cap1)
	cameras = append(cameras, &cap1)

	tick := snippets.Tick{}
	tick.SetUp(200000)
	tb.AddSource("tick", &tick)
	ticker = &tick

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
	itr := snippets.Integrate{
		ConfigPath: confPath + "integrate[0].json",
	}
	tb.AddBox("integrate", &itr).Input("recognize_caffe")

	sender := snippets.DataSender{}
	sender.SetUp(confPath + "integrate[0].json")
	tb.AddSink("data_sender", &sender).Input("integrate")

	return tb.Build()
}

// attention permanent loop
func memStats(m *runtime.MemStats, ctx *core.Context) {
	for {
		select {
		case <-time.After(time.Second):
			runtime.ReadMemStats(m)
			ctx.Logger.Log(core.Debug, "memstats,%d,%d,%d,%d", m.HeapSys, m.HeapAlloc, m.HeapIdle, m.HeapReleased)
		}
	}
}

func tickerStopper() {
stopper:
	for {
		select {
		case <-time.After(time.Second):
			for _, cap := range cameras {
				if !cap.IsStopped() {
					continue stopper
				}
			}
			ticker.ForcedStop()
			return
		}
	}
}

func main() {
	// go does not optimize goroutine threads with CPU core
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus * 2)

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

	// performance
	//var m runtime.MemStats
	//go memStats(&m, &ctx)

	go tickerStopper()

	topoloby.Run(&ctx)
}
