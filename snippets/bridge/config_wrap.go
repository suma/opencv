package bridge

/*
#include "config_bridge.h"
*/
import "C"

type FrameProcessorConfig struct {
	p C.FrameProcessorConfig
}

func FrameProcessorConfig_New(config string) FrameProcessorConfig {
	return FrameProcessorConfig{
		p: C.FrameProcessorConfig_New(C.CString(config)),
	}
}

func (fpc *FrameProcessorConfig) Delete() {
	C.FrameProcessorConfig_Delete(fpc.p)
}
