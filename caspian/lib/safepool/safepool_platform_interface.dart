import 'package:plugin_platform_interface/plugin_platform_interface.dart';

import 'safepool_method_channel.dart';

abstract class SafepoolPlatform extends PlatformInterface {
  /// Constructs a SafepoolPlatform.
  SafepoolPlatform() : super(token: _token);

  static final Object _token = Object();

  static SafepoolPlatform _instance = MethodChannelSafepool();

  /// The default instance of [SafepoolPlatform] to use.
  ///
  /// Defaults to [MethodChannelSafepool].
  static SafepoolPlatform get instance => _instance;

  /// Platform-specific implementations should set this with their own
  /// platform-specific class that extends [SafepoolPlatform] when
  /// they register themselves.
  static set instance(SafepoolPlatform instance) {
    PlatformInterface.verifyToken(instance, _token);
    _instance = instance;
  }

  Future<String?> getPlatformVersion() {
    throw UnimplementedError('platformVersion() has not been implemented.');
  }
}
