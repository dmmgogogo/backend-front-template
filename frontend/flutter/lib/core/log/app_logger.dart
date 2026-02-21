import 'package:flutter/foundation.dart';
import 'package:logger/logger.dart';

/// Unified logger entry for the whole app.
///
/// - debug/profile: output all levels
/// - release: output warning and error only
class AppLogger {
  AppLogger._();

  static const bool _verbose = bool.fromEnvironment(
    'VERBOSE_P2P_LOGS',
    defaultValue: false,
  );

  static final Logger _logger = Logger(
    printer: SimplePrinter(colors: false, printTime: true),
    level: kReleaseMode ? Level.warning : (_verbose ? Level.trace : Level.info),
  );
  static const int _maxBufferedLogs = 800;
  static final List<String> _bufferedLogs = <String>[];
  static final ValueNotifier<int> _logRevision = ValueNotifier<int>(0);
  static final ValueNotifier<bool> _captureEnabled = ValueNotifier<bool>(false);

  static LogScope scope(String name) => LogScope._(name);
  static ValueListenable<int> get logRevision => _logRevision;
  static ValueListenable<bool> get captureEnabledListenable => _captureEnabled;
  static bool get isCaptureEnabled => _captureEnabled.value;
  static List<String> get bufferedLogs =>
      List<String>.unmodifiable(_bufferedLogs);

  static void setCaptureEnabled(bool enabled) {
    if (_captureEnabled.value == enabled) return;
    _captureEnabled.value = enabled;
    i('AppLogger', enabled ? 'capture.enabled' : 'capture.disabled');
  }

  static void clearBufferedLogs() {
    _bufferedLogs.clear();
    _logRevision.value++;
  }

  static void t(
    String scope,
    String event, {
    Map<String, Object?> ctx = const {},
  }) {
    final content = _format(scope, event, ctx);
    _logger.t(content);
    _capture('T', content);
  }

  static void d(
    String scope,
    String event, {
    Map<String, Object?> ctx = const {},
  }) {
    final content = _format(scope, event, ctx);
    _logger.d(content);
    _capture('D', content);
  }

  static void i(
    String scope,
    String event, {
    Map<String, Object?> ctx = const {},
  }) {
    final content = _format(scope, event, ctx);
    _logger.i(content);
    _capture('I', content);
  }

  static void w(
    String scope,
    String event, {
    Object? error,
    StackTrace? stackTrace,
    Map<String, Object?> ctx = const {},
  }) {
    final content = _format(scope, event, ctx);
    _logger.w(content, error: error, stackTrace: stackTrace);
    _capture('W', '$content ${error ?? ''}'.trim());
  }

  static void e(
    String scope,
    String event, {
    Object? error,
    StackTrace? stackTrace,
    Map<String, Object?> ctx = const {},
  }) {
    final content = _format(scope, event, ctx);
    _logger.e(content, error: error, stackTrace: stackTrace);
    _capture('E', '$content ${error ?? ''}'.trim());
  }

  static String redactSecret(String secret) {
    if (secret.isEmpty) return '';
    if (secret.length <= 10) return '*' * secret.length;
    return '${secret.substring(0, 6)}...${secret.substring(secret.length - 4)}';
  }

  static String shortPeer(String peerId) {
    if (peerId.length <= 12) return peerId;
    return '${peerId.substring(0, 8)}...${peerId.substring(peerId.length - 4)}';
  }

  static String _format(String scope, String event, Map<String, Object?> ctx) {
    if (ctx.isEmpty) {
      return '[$scope] $event';
    }
    final entries = ctx.entries
        .map((entry) => '${entry.key}=${entry.value}')
        .join(' ');
    return '[$scope] $event | $entries';
  }

  static void _capture(String level, String content) {
    if (!_captureEnabled.value) return;
    final line = '[${DateTime.now().toIso8601String()}][$level] $content';
    _bufferedLogs.add(line);
    if (_bufferedLogs.length > _maxBufferedLogs) {
      _bufferedLogs.removeRange(0, _bufferedLogs.length - _maxBufferedLogs);
    }
    _logRevision.value++;
  }
}

class LogScope {
  LogScope._(this._scope);

  final String _scope;

  void t(String event, {Map<String, Object?> ctx = const {}}) =>
      AppLogger.t(_scope, event, ctx: ctx);

  void d(String event, {Map<String, Object?> ctx = const {}}) =>
      AppLogger.d(_scope, event, ctx: ctx);

  void i(String event, {Map<String, Object?> ctx = const {}}) =>
      AppLogger.i(_scope, event, ctx: ctx);

  void w(
    String event, {
    Object? error,
    StackTrace? stackTrace,
    Map<String, Object?> ctx = const {},
  }) => AppLogger.w(
    _scope,
    event,
    error: error,
    stackTrace: stackTrace,
    ctx: ctx,
  );

  void e(
    String event, {
    Object? error,
    StackTrace? stackTrace,
    Map<String, Object?> ctx = const {},
  }) => AppLogger.e(
    _scope,
    event,
    error: error,
    stackTrace: stackTrace,
    ctx: ctx,
  );
}
