import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart';

import 'safepool_platform_interface.dart';

/// An implementation of [SafepoolPlatform] that uses method channels.
class MethodChannelSafepool extends SafepoolPlatform {
  /// The method channel used to interact with the native platform.
  @visibleForTesting
  final methodChannel = const MethodChannel('safepool');

  @override
  Future<String?> getPlatformVersion() async {
    final version = await methodChannel.invokeMethod<String>('getPlatformVersion');
    return version;
  }
}
