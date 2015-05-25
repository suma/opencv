package bridge

/*
#cgo LDFLAGS: -ljsonconfig
#include "config_bridge.h"
*/
import "C"

type FrameProcessorConfig struct {
	p C.FrameProcessorConfig
}

type DetectorConfig struct {
	p C.DetectorConfig
}

type RecognizeConfigTaggers struct {
	p C.RecognizeConfigTaggers
}

type IntegratorConfig struct {
	p C.IntegratorConfig
}

func FrameProcessorConfig_New(config string) FrameProcessorConfig {
	return FrameProcessorConfig{
		p: C.FrameProcessorConfig_New(C.CString(config)),
	}
}

func (fpc *FrameProcessorConfig) Delete() {
	C.FrameProcessorConfig_Delete(fpc.p)
}

func DetectorConfig_New(config string) DetectorConfig {
	return DetectorConfig{
		p: C.DetectorConfig_New(C.CString(config)),
	}
}

func (dc *DetectorConfig) Delete() {
	C.DetectorConfig_Delete(dc.p)
}

func RecognizeConfigTaggers_New(config string) RecognizeConfigTaggers {
	return RecognizeConfigTaggers{
		p: C.RecognizeConfigTaggers_New(C.CString(config)),
	}
}

func (ta *RecognizeConfigTaggers) Delete() {
	C.RecognizeConfigTaggers_Delete(ta.p)
}

func IntegratorConfig_New(config string) IntegratorConfig {
	return IntegratorConfig{
		p: C.IntegratorConfig_New(C.CString(config)),
	}
}

func (ic *IntegratorConfig) Delete() {
	C.IntegratorConfig_Delete(ic.p)
}
