// ========== iOS 内购支付代码（复制到 lib/features/support/support_service.dart）==========
// 后端接口: POST /api/backend/support/ios/verify（见 backend CLAUDE.md「iOS 支付」）
// 依赖: in_app_purchase；iOS 需 in_app_purchase_storekit
//
// 1. 在 pubspec.yaml 添加:
//    dependencies:
//      in_app_purchase: ^3.2.0
//
// 2. 商品 ID 与后端 controllers/backend/support_ios.go 中 iosProductAmount 一致
// 3. 请求需带 token（与后端 JWT 一致），HttpClient 的 AuthInterceptor 会自动携带
//
// 使用示例:
//   await SupportService().purchaseIOSSupport(amount: 10, onDebugEvent: (msg) => debugPrint(msg));
//

import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'package:in_app_purchase/in_app_purchase.dart';
import 'package:in_app_purchase_storekit/in_app_purchase_storekit.dart';
import '../../core/network/http_client.dart';

class SupportPurchaseException implements Exception {
  final String stage;
  final String message;
  const SupportPurchaseException(this.stage, this.message);
  @override
  String toString() => '[$stage] $message';
}

class SupportService {
  final InAppPurchase _iap = InAppPurchase.instance;

  static const Map<String, double> _iosProductAmount = {
    'com.yourapp.support.1': 1,
    'com.yourapp.support.5': 5,
    'com.yourapp.support.10': 10,
    'com.yourapp.support.50': 50,
    'com.yourapp.support.100': 100,
    'com.yourapp.support.300': 300,
    'com.yourapp.support.500': 500,
    'com.yourapp.support.1000': 1000,
  };

  Future<void> purchaseIOSSupport({
    required double amount,
    void Function(String message)? onDebugEvent,
  }) async {
    final watch = Stopwatch()..start();
    void logStep(String message) {
      onDebugEvent?.call('[+${watch.elapsedMilliseconds}ms] $message');
    }

    logStep('start amount=\$${amount.toStringAsFixed(2)}');
    final available = await _iap.isAvailable().timeout(const Duration(seconds: 12));
    if (!available) throw const SupportPurchaseException('availability', 'App Store payment is unavailable');

    final targetProductId = _pickNearestProduct(amount);
    logStep('target_product_id=$targetProductId');
    final response = await _iap.queryProductDetails(_iosProductAmount.keys.toSet()).timeout(const Duration(seconds: 15));
    if (response.error != null) throw SupportPurchaseException('product_query', response.error!.message);

    final product = response.productDetails.where((p) => p.id == targetProductId).cast<ProductDetails?>().firstWhere((p) => p != null, orElse: () => null);
    if (product == null) throw SupportPurchaseException('product_query', 'Target product not found: $targetProductId');
    logStep('product_found id=${product.id}');

    final completer = Completer<void>();
    late final StreamSubscription<List<PurchaseDetails>> sub;
    var hasRestoredEvent = false;

    Future<void> finishWithError(Object e) async {
      if (!completer.isCompleted) completer.completeError(e);
      await sub.cancel();
    }

    Future<void> finishSuccess() async {
      if (!completer.isCompleted) completer.complete();
      await sub.cancel();
    }

    sub = _iap.purchaseStream.listen((purchases) async {
      for (final purchase in purchases) {
        if (purchase.productID != targetProductId) continue;
        switch (purchase.status) {
          case PurchaseStatus.pending:
            break;
          case PurchaseStatus.error:
            await finishWithError(SupportPurchaseException('purchase_stream', purchase.error?.message ?? 'Purchase failed'));
            break;
          case PurchaseStatus.canceled:
            await finishWithError(const SupportPurchaseException('purchase_stream', 'Purchase canceled'));
            break;
          case PurchaseStatus.purchased:
          case PurchaseStatus.restored:
            if (purchase.status == PurchaseStatus.restored) hasRestoredEvent = true;
            try {
              await _verifyIOSPurchase(purchase, onDebugEvent: onDebugEvent);
              await finishSuccess();
            } catch (e) {
              await finishWithError(e);
            } finally {
              if (purchase.pendingCompletePurchase) await _iap.completePurchase(purchase);
            }
            break;
        }
      }
    });

    final launched = await _iap.buyConsumable(purchaseParam: PurchaseParam(productDetails: product));
    if (!launched) {
      await sub.cancel();
      throw const SupportPurchaseException('buy_launch', 'Unable to launch purchase flow');
    }

    return completer.future.timeout(
      const Duration(seconds: 30),
      onTimeout: () async {
        await sub.cancel();
        throw SupportPurchaseException('timeout', hasRestoredEvent ? 'Only restored received. Clear Sandbox and retry.' : 'Purchase status timeout.');
      },
    );
  }

  String _pickNearestProduct(double amount) {
    final sorted = _iosProductAmount.entries.toList()..sort((a, b) => a.value.compareTo(b.value));
    for (final e in sorted) {
      if (amount <= e.value) return e.key;
    }
    return sorted.last.key;
  }

  Future<void> _verifyIOSPurchase(PurchaseDetails purchase, {void Function(String message)? onDebugEvent}) async {
    String receipt = purchase.verificationData.localVerificationData.trim();
    if (receipt.isEmpty) receipt = purchase.verificationData.serverVerificationData.trim();
    if (receipt.isEmpty) throw const SupportPurchaseException('receipt', 'Missing receipt data');

    if (_looksLikeJws(receipt) || _looksLikeJson(receipt)) {
      if (Platform.isIOS) {
        final addition = _iap.getPlatformAddition<InAppPurchaseStoreKitPlatformAddition>();
        final refreshed = await addition.refreshPurchaseVerificationData();
        final refreshedReceipt = refreshed?.localVerificationData.trim();
        if (refreshedReceipt != null && refreshedReceipt.isNotEmpty) receipt = refreshedReceipt;
      }
      if (_looksLikeJws(receipt) || _looksLikeJson(receipt)) throw const SupportPurchaseException('receipt', 'Unsupported receipt format, need base64 app receipt');
    }

    final transactionID = (purchase.purchaseID ?? '').isNotEmpty ? purchase.purchaseID! : '${purchase.productID}_${base64Url.encode(utf8.encode(receipt.length > 32 ? receipt.substring(0, 32) : receipt))}';
    onDebugEvent?.call('verify_request product=${purchase.productID} tx=$transactionID');

    final response = await HttpClient().dio.post('/support/ios/verify', data: {
      'product_id': purchase.productID,
      'transaction_id': transactionID,
      'receipt_data': receipt,
    });

    final data = response.data is Map ? response.data as Map<String, dynamic> : null;
    final code = data?['code'];
    onDebugEvent?.call('verify_response code=$code msg=${data?['msg']}');
    if (code != 200) throw SupportPurchaseException('receipt_verify', data?['msg']?.toString() ?? 'Verification failed');
  }

  bool _looksLikeJws(String v) => v.isNotEmpty && v.split('.').length == 3;
  bool _looksLikeJson(String v) {
    final t = v.trimLeft();
    return t.startsWith('{') || t.startsWith('[');
  }
}
