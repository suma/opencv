#pragma once

#include "util.h"
#include <cstdlib>
#include <msgpack.hpp>

ByteArray toByteArray(const char* buf, int len);

template<class T>
ByteArray serializeObject(const T& obj) {
  msgpack::sbuffer buf;
  msgpack::packer<msgpack::sbuffer> pk(&buf);
  pk.pack(*obj);
  return toByteArray(buf.data(), buf.size());
}

template<class T>
T* deserializeObject(const ByteArray& src) {
  msgpack::unpacked msg;
  msgpack::unpack(&msg, src.data, src.length);
  msgpack::object obj = msg.get();
  T* ret = new T();
  obj.convert(ret);
  return ret;
}
