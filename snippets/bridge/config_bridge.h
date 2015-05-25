#ifndef _CONFIG_BRIDGE_H_
#define _CONFIG_BRIDGE_H_

#ifdef __cplusplus
extern "C" {
#endif

typedef void* FrameProcessorConfig;
typedef void* DetectorConfig;
typedef void* RecognizeConfigTaggers;
typedef void* IntegratorConfig;

FrameProcessorConfig FrameProcessorConfig_New(const char *config);
void FrameProcessorConfig_Delete(FrameProcessorConfig config);

DetectorConfig DetectorConfig_New(const char *config);
void DetectorConfig_Delete(DetectorConfig config);

RecognizeConfigTaggers RecognizeConfigTaggers_New(const char *config);
void RecognizeConfigTaggers_Delete(RecognizeConfigTaggers taggers);

IntegratorConfig IntegratorConfig_New(const char *config);
void IntegratorConfig_Delete(IntegratorConfig config);

#ifdef __cplusplus
}
#endif

#endif //_CONFIG_BRIDGE_H_