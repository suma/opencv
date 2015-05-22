#ifndef _CONFIG_BRIDGE_H_
#define _CONFIG_BRIDGE_H_

#ifdef __cplusplus
extern "C" {
#endif

typedef void* FrameProcessorConfig;

FrameProcessorConfig FrameProcessorConfig_New(const char *config);
void FrameProcessorConfig_Delete(FrameProcessorConfig config);

#ifdef __cplusplus
}
#endif

#endif //_CONFIG_BRIDGE_H_