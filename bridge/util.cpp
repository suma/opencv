#include "util.hpp"
#include <cstdlib>

void ByteArray_Release(ByteArray buf) {
  delete[] buf.data;
}

ByteArray toByteArray(const char* buf, int len) {
  ByteArray ret = {new char[len], len};
  memcpy(ret.data, buf, len);
  return ret;
}
