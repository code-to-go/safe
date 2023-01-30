import 'dart:ffi' as ffi; // For FFI
import 'dart:io'; // For Platform.isX
import 'dart:convert';
import "package:ffi/ffi.dart";
import 'safepool_platform_interface.dart';

final ffi.DynamicLibrary lib = getLibrary();

ffi.DynamicLibrary getLibrary() {
  if (Platform.isAndroid) {
    return ffi.DynamicLibrary.open('libsafepool.so');
  }
  if (Platform.isLinux) {
    return ffi.DynamicLibrary.open('linux/libs/amd64/libsafepool.so');
  }
  return ffi.DynamicLibrary.process();
}

class Safepool {
  Future<String?> getPlatformVersion() {
    return SafepoolPlatform.instance.getPlatformVersion();
  }
}

class CResult extends ffi.Struct {
  external ffi.Pointer<Utf8> res;
  external ffi.Pointer<Utf8> err;

  void unwrapVoid() {
    if (err.address != 0) {
      throw CException(err.toDartString());
    }
  }

  String unwrapString() {
    unwrapVoid();
    if (res.address == 0) {
      return "";
    }

    return jsonDecode(res.toDartString()) as String;
  }

  Map<String, dynamic> unwrapMap() {
    unwrapVoid();
    if (res.address == 0) {
      return {};
    }

    return jsonDecode(res.toDartString()) as Map<String, dynamic>;
  }

  List<dynamic> unwrapList() {
    if (err.address != 0) {
      throw CException(err.toDartString());
    }
    if (res.address == 0) {
      return [];
    }

    var ls = jsonDecode(res.toDartString());

    return ls == null ? [] : ls as List<dynamic>;
  }
}

class CException implements Exception {
  String msg;
  CException(this.msg);
}

typedef Start = CResult Function(ffi.Pointer<Utf8>);
void start(String dbPath) {
  var startC = lib.lookupFunction<Start, Start>("start");
  startC(dbPath.toNativeUtf8()).unwrapVoid();
}

typedef GetSelfId = CResult Function();
String getSelfId() {
  var getSelfIdC = lib.lookupFunction<GetSelfId, GetSelfId>("getSelfId");
  return getSelfIdC().unwrapString();
}

typedef GetPoolList = CResult Function();
List<String> getPoolList() {
  var getPoolListC = lib.lookupFunction<GetSelfId, GetSelfId>("getPoolList");
  return getPoolListC().unwrapList().map((e) => e as String).toList();
}

typedef CreatePool = CResult Function(ffi.Pointer<Utf8>);
List<String> createPool() {
  var getPoolListC = lib.lookupFunction<GetSelfId, GetSelfId>("getPoolList");
  return getPoolListC().unwrapList().map((e) => e as String).toList();
}
