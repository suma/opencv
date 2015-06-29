#ifndef _BRIDGE_UTIL_HPP_
#define _BRIDGE_UTIL_HPP_

#include "util.h"
#include <cstdlib>
#include <msgpack.hpp>
#include <pficommon/text/json.h>
#include <jsonconfig.hpp>

ByteArray toByteArray(const char* buf, int len);

template<class T>
ByteArray serializeObject(const T& obj) {
  msgpack::sbuffer buf;
  msgpack::packer<msgpack::sbuffer> pk(&buf);
  pk.pack(obj);
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

template <class Type>
Type load_json(const char *config) {
  std::stringstream ss(config);
  pfi::text::json::json config_raw;
  ss >> config_raw;
  return jsonconfig::config_cast<Type>(jsonconfig::config_root(config_raw));
}

#endif // _BRIDGE_UTIL_HPP_
