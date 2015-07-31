#ifndef _BRIDGE_UTIL_H_
#define _BRIDGE_UTIL_H_

#ifdef __cplusplus
extern "C" {
#endif

typedef struct String {
  const char* str;
  int length;
} String;
struct ByteArray{
  char *data;
  int length;
};

void ByteArray_Release(struct ByteArray buf);

#ifdef __cplusplus
}
#endif

#endif //_BRIDGE_UTIL_H_
