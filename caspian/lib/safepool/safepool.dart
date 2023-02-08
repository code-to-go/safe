import 'dart:ffi' as ffi; // For FFI
import 'dart:io'; // For Platform.isX
import 'dart:convert';
import 'dart:typed_data';
import 'package:caspian/safepool/safepool_def.dart';
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

  int unwrapInt() {
    unwrapVoid();

    return jsonDecode(res.toDartString()) as int;
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

  @override
  String toString() {
    return msg;
  }
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
  var getPoolListC =
      lib.lookupFunction<GetPoolList, GetPoolList>("getPoolList");
  return getPoolListC().unwrapList().map((e) => e as String).toList();
}

typedef CreatePool = CResult Function(ffi.Pointer<Utf8>, ffi.Pointer<Utf8>);
void createPool(PoolConfig c, List<String> apps) {
  var createPoolC = lib.lookupFunction<CreatePool, CreatePool>("createPool");
  var appsJson = "[${apps.map((e) => jsonEncode(e)).toList().join(",")}]";
  var cJson = jsonEncode(c.toJson());
  return createPoolC(cJson.toNativeUtf8(), appsJson.toNativeUtf8())
      .unwrapVoid();
}

typedef AddPool = CResult Function(ffi.Pointer<Utf8>);
PoolConfig addPool(String token) {
  var addPoolC = lib.lookupFunction<AddPool, AddPool>("addPool");
  var m = addPoolC(token.toNativeUtf8()).unwrapMap();
  return PoolConfig.fromJson(m);
}

typedef GetPool = CResult Function(ffi.Pointer<Utf8>);
Pool getPool(String name) {
  var getPoolC = lib.lookupFunction<GetPool, GetPool>("getPool");
  var m = getPoolC(name.toNativeUtf8()).unwrapMap();
  return Pool.fromJson(m);
}

typedef GetMessagesC = CResult Function(
    ffi.Pointer<Utf8>, ffi.Int64, ffi.Int64, ffi.Int32);
typedef GetMessages = CResult Function(ffi.Pointer<Utf8>, int, int, int);
List<Message> getMessages(
    String poolName, int afterId, int beforeId, int limit) {
  var getMessagesC =
      lib.lookupFunction<GetMessagesC, GetMessages>("getMessages");
  var m = getMessagesC(poolName.toNativeUtf8(), afterId, beforeId, limit)
      .unwrapList();

  return m.map((e) => Message.fromJson(e as Map<String, dynamic>)).toList();
}

typedef PostMessage = CResult Function(
    ffi.Pointer<Utf8>, ffi.Pointer<Utf8>, ffi.Pointer<Utf8>, ffi.Pointer<Utf8>);
int postMessage(
    String poolName, String contentType, String text, Uint8List binary) {
  var postMessageC =
      lib.lookupFunction<PostMessage, PostMessage>("postMessage");

  var r = postMessageC(
    poolName.toNativeUtf8(),
    contentType.toNativeUtf8(),
    text.toNativeUtf8(),
    base64Encode(binary).toNativeUtf8(),
  );
  return r.unwrapInt();
}
