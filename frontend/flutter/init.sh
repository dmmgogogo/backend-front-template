#!/bin/bash
#
# Flutter 项目一键初始化脚本
# 用法: bash init.sh <project_name> [org_name]
# 示例: bash init.sh my_app com.yourcompany
#

set -e

# ========== 参数 ==========
PROJECT_NAME="${1:-my_app}"
ORG_NAME="${2:-com.yourcompany}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# ========== 颜色 ==========
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
warn()    { echo -e "${YELLOW}[WARN]${NC} $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# ========== 环境检查 ==========
info "检查 Flutter 环境..."
if ! command -v flutter &> /dev/null; then
    error "未安装 Flutter，请先安装: https://flutter.dev/docs/get-started/install"
fi

FLUTTER_VERSION=$(flutter --version 2>&1 | head -1)
success "Flutter 已安装: $FLUTTER_VERSION"

# ========== 确认参数 ==========
echo ""
echo "=================================="
echo "  Flutter 项目初始化"
echo "=================================="
echo "  项目名称: $PROJECT_NAME"
echo "  组织名称: $ORG_NAME"
echo "  目标目录: $SCRIPT_DIR"
echo "=================================="
echo ""
read -p "确认创建？(y/N) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    info "已取消"
    exit 0
fi

# ========== 创建 Flutter 项目 ==========
info "创建 Flutter 项目..."

# 在临时目录创建，然后移动内容到当前目录
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"
flutter create --org "$ORG_NAME" --project-name "$PROJECT_NAME" "$PROJECT_NAME" --platforms=ios,android

# 移动内容到脚本所在目录（保留 CLAUDE.md 和 init.sh）
info "移动项目文件..."
cd "$TEMP_DIR/$PROJECT_NAME"
for item in *; do
    if [ "$item" != "CLAUDE.md" ] && [ "$item" != "init.sh" ]; then
        cp -r "$item" "$SCRIPT_DIR/"
    fi
done
# 移动隐藏文件（.gitignore 等）
for item in .[!.]*; do
    [ -e "$item" ] && cp -r "$item" "$SCRIPT_DIR/"
done

# 清理临时目录
rm -rf "$TEMP_DIR"
cd "$SCRIPT_DIR"

success "Flutter 项目创建完成"

# ========== 创建目录结构 ==========
info "创建标准目录结构..."

# config/
mkdir -p lib/config
# core/
mkdir -p lib/core/network/interceptors
mkdir -p lib/core/storage
mkdir -p lib/core/utils
mkdir -p lib/core/widgets
# features/
mkdir -p lib/features/auth/{data,models,providers,pages}
mkdir -p lib/features/home/{data,models,providers,pages}
mkdir -p lib/features/profile/{data,models,providers,pages}
# l10n/
mkdir -p lib/l10n
# assets/
mkdir -p assets/images
mkdir -p assets/icons
mkdir -p assets/fonts

success "目录结构创建完成"

# ========== 创建基础文件 ==========
info "创建基础文件..."

# --- config/env.dart ---
cat > lib/config/env.dart << 'DART_EOF'
/// 环境配置
enum Env { dev, staging, prod }

class EnvConfig {
  static Env current = Env.dev;

  static String get baseUrl {
    switch (current) {
      case Env.dev:
        return 'http://localhost:8282/api/backend';
      case Env.staging:
        return 'https://staging-api.yourcompany.com/api/backend';
      case Env.prod:
        return 'https://api.yourcompany.com/api/backend';
    }
  }

  /// 通用接口地址（文件上传等）
  static String get commonBaseUrl {
    switch (current) {
      case Env.dev:
        return 'http://localhost:8282/api/common';
      case Env.staging:
        return 'https://staging-api.yourcompany.com/api/common';
      case Env.prod:
        return 'https://api.yourcompany.com/api/common';
    }
  }

  static bool get isDebug => current == Env.dev;
}
DART_EOF

# --- config/theme.dart ---
cat > lib/config/theme.dart << 'DART_EOF'
import 'package:flutter/material.dart';

/// 应用主题配置
class AppTheme {
  static const Color primaryColor = Color(0xFF1890FF);
  static const Color errorColor = Color(0xFFFF4D4F);
  static const Color successColor = Color(0xFF52C41A);
  static const Color warningColor = Color(0xFFFAAD14);

  static ThemeData get lightTheme {
    return ThemeData(
      primaryColor: primaryColor,
      colorScheme: ColorScheme.fromSeed(
        seedColor: primaryColor,
        brightness: Brightness.light,
      ),
      useMaterial3: true,
    );
  }

  static ThemeData get darkTheme {
    return ThemeData(
      primaryColor: primaryColor,
      colorScheme: ColorScheme.fromSeed(
        seedColor: primaryColor,
        brightness: Brightness.dark,
      ),
      useMaterial3: true,
    );
  }
}
DART_EOF

# --- config/routes.dart ---
cat > lib/config/routes.dart << 'DART_EOF'
import 'package:go_router/go_router.dart';
import '../features/auth/pages/login_page.dart';
import '../features/home/pages/home_page.dart';

/// 路由配置
final GoRouter appRouter = GoRouter(
  initialLocation: '/home',
  routes: [
    GoRoute(
      path: '/login',
      builder: (context, state) => const LoginPage(),
    ),
    GoRoute(
      path: '/home',
      builder: (context, state) => const HomePage(),
    ),
  ],
);
DART_EOF

# --- core/network/api_response.dart ---
cat > lib/core/network/api_response.dart << 'DART_EOF'
/// 统一 API 响应模型
/// 与后端返回格式对齐: { "code": 200, "msg": "success", "data": {...} }
class ApiResponse<T> {
  final int code;
  final String msg;
  final T? data;

  ApiResponse({required this.code, required this.msg, this.data});

  bool get isSuccess => code == 200;

  factory ApiResponse.fromJson(
    Map<String, dynamic> json,
    T Function(dynamic)? fromData,
  ) {
    return ApiResponse(
      code: json['code'] ?? -1,
      msg: json['msg'] ?? '',
      data: json['data'] != null && fromData != null
          ? fromData(json['data'])
          : json['data'],
    );
  }
}
DART_EOF

# --- core/network/http_client.dart ---
cat > lib/core/network/http_client.dart << 'DART_EOF'
import 'package:dio/dio.dart';
import '../../config/env.dart';
import 'interceptors/auth_interceptor.dart';
import 'interceptors/error_interceptor.dart';

/// Dio HTTP 客户端封装
class HttpClient {
  static HttpClient? _instance;
  late Dio dio;

  HttpClient._internal() {
    dio = Dio(BaseOptions(
      baseUrl: EnvConfig.baseUrl,
      connectTimeout: const Duration(seconds: 15),
      receiveTimeout: const Duration(seconds: 15),
      headers: {
        'Content-Type': 'application/json',
      },
    ));

    // 添加拦截器
    dio.interceptors.addAll([
      AuthInterceptor(),
      ErrorInterceptor(),
      if (EnvConfig.isDebug) LogInterceptor(requestBody: true, responseBody: true),
    ]);
  }

  factory HttpClient() {
    _instance ??= HttpClient._internal();
    return _instance!;
  }

  /// GET 请求
  Future<Response> get(String path, {Map<String, dynamic>? params}) {
    return dio.get(path, queryParameters: params);
  }

  /// POST 请求
  Future<Response> post(String path, {dynamic data}) {
    return dio.post(path, data: data);
  }

  /// 上传文件
  Future<Response> upload(String path, FormData formData) {
    return dio.post(path, data: formData);
  }
}
DART_EOF

# --- core/network/interceptors/auth_interceptor.dart ---
cat > lib/core/network/interceptors/auth_interceptor.dart << 'DART_EOF'
import 'package:dio/dio.dart';
import '../../storage/sp_util.dart';

/// 认证拦截器 - 自动注入 token
class AuthInterceptor extends Interceptor {
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    final token = SpUtil.getToken();
    if (token != null && token.isNotEmpty) {
      options.headers['token'] = token;
    }

    // 注入语言标识
    final lang = SpUtil.getLanguage() ?? 'zh';
    options.headers['Accept-Language'] = lang;

    handler.next(options);
  }
}
DART_EOF

# --- core/network/interceptors/error_interceptor.dart ---
cat > lib/core/network/interceptors/error_interceptor.dart << 'DART_EOF'
import 'package:dio/dio.dart';
import '../../storage/sp_util.dart';

/// 错误拦截器 - 统一错误处理
class ErrorInterceptor extends Interceptor {
  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    final data = response.data;
    if (data is Map<String, dynamic>) {
      final code = data['code'];
      if (code == 401) {
        // Token 过期，清除登录态
        SpUtil.clearToken();
        // TODO: 跳转登录页
      }
    }
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    String message;
    switch (err.type) {
      case DioExceptionType.connectionTimeout:
        message = '连接超时';
        break;
      case DioExceptionType.receiveTimeout:
        message = '响应超时';
        break;
      case DioExceptionType.connectionError:
        message = '网络连接异常';
        break;
      default:
        message = '网络异常: ${err.message}';
    }
    // TODO: 显示 toast 提示
    handler.next(err);
  }
}
DART_EOF

# --- core/storage/sp_util.dart ---
cat > lib/core/storage/sp_util.dart << 'DART_EOF'
import 'package:shared_preferences/shared_preferences.dart';

/// SharedPreferences 封装
class SpUtil {
  static SharedPreferences? _prefs;

  static Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
  }

  // Token
  static const String _keyToken = 'token';
  static String? getToken() => _prefs?.getString(_keyToken);
  static Future<bool> setToken(String token) =>
      _prefs!.setString(_keyToken, token);
  static Future<bool> clearToken() => _prefs!.remove(_keyToken);

  // Language
  static const String _keyLanguage = 'language';
  static String? getLanguage() => _prefs?.getString(_keyLanguage);
  static Future<bool> setLanguage(String lang) =>
      _prefs!.setString(_keyLanguage, lang);

  // 通用
  static Future<bool> clear() => _prefs!.clear();
}
DART_EOF

# --- core/widgets/loading.dart ---
cat > lib/core/widgets/loading.dart << 'DART_EOF'
import 'package:flutter/material.dart';

/// 通用 Loading 组件
class LoadingWidget extends StatelessWidget {
  final String? message;

  const LoadingWidget({super.key, this.message});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const CircularProgressIndicator(),
          if (message != null) ...[
            const SizedBox(height: 16),
            Text(message!, style: const TextStyle(color: Colors.grey)),
          ],
        ],
      ),
    );
  }
}
DART_EOF

# --- core/widgets/empty_view.dart ---
cat > lib/core/widgets/empty_view.dart << 'DART_EOF'
import 'package:flutter/material.dart';

/// 空状态组件
class EmptyView extends StatelessWidget {
  final String message;
  final IconData icon;

  const EmptyView({
    super.key,
    this.message = '暂无数据',
    this.icon = Icons.inbox_outlined,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(icon, size: 64, color: Colors.grey[400]),
          const SizedBox(height: 16),
          Text(message, style: TextStyle(color: Colors.grey[600], fontSize: 16)),
        ],
      ),
    );
  }
}
DART_EOF

# --- features/auth/pages/login_page.dart ---
cat > lib/features/auth/pages/login_page.dart << 'DART_EOF'
import 'package:flutter/material.dart';

/// 登录页（骨架）
class LoginPage extends StatelessWidget {
  const LoginPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('登录')),
      body: const Center(child: Text('Login Page - TODO')),
    );
  }
}
DART_EOF

# --- features/home/pages/home_page.dart ---
cat > lib/features/home/pages/home_page.dart << 'DART_EOF'
import 'package:flutter/material.dart';

/// 首页（骨架）
class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('首页')),
      body: const Center(child: Text('Home Page - TODO')),
    );
  }
}
DART_EOF

# --- app.dart ---
cat > lib/app.dart << 'DART_EOF'
import 'package:flutter/material.dart';
import 'config/theme.dart';
import 'config/routes.dart';

/// App 根组件
class App extends StatelessWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'My App',
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.system,
      routerConfig: appRouter,
      debugShowCheckedModeBanner: false,
    );
  }
}
DART_EOF

# --- 覆盖 main.dart ---
cat > lib/main.dart << 'DART_EOF'
import 'package:flutter/material.dart';
import 'app.dart';
import 'config/env.dart';
import 'core/storage/sp_util.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // 环境配置
  EnvConfig.current = Env.dev;

  // 初始化本地存储
  await SpUtil.init();

  runApp(const App());
}
DART_EOF

# --- l10n ---
cat > lib/l10n/intl_zh.arb << 'JSON_EOF'
{
  "@@locale": "zh",
  "appTitle": "我的应用",
  "login": "登录",
  "logout": "退出登录",
  "home": "首页",
  "profile": "个人中心"
}
JSON_EOF

cat > lib/l10n/intl_en.arb << 'JSON_EOF'
{
  "@@locale": "en",
  "appTitle": "My App",
  "login": "Login",
  "logout": "Logout",
  "home": "Home",
  "profile": "Profile"
}
JSON_EOF

success "基础文件创建完成"

# ========== 更新 pubspec.yaml ==========
info "更新 pubspec.yaml 依赖..."

# 使用 flutter pub add 添加依赖（自动获取最新版本）
cd "$SCRIPT_DIR"
flutter pub add dio provider go_router shared_preferences json_annotation logger 2>/dev/null || warn "部分依赖添加失败，请手动检查"
flutter pub add --dev json_serializable build_runner 2>/dev/null || warn "部分 dev 依赖添加失败，请手动检查"

# 添加 assets 到 pubspec.yaml（如果还没有）
if ! grep -q "assets/images/" pubspec.yaml 2>/dev/null; then
    # 在 flutter: 块中添加 assets
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' '/^flutter:/,/^[^ ]/ {
            /uses-material-design: true/a\
\  assets:\
\    - assets/images/\
\    - assets/icons/
        }' pubspec.yaml 2>/dev/null || warn "请手动在 pubspec.yaml 的 flutter 块中添加 assets 配置"
    else
        sed -i '/^flutter:/,/^[^ ]/ {
            /uses-material-design: true/a\  assets:\n    - assets/images/\n    - assets/icons/
        }' pubspec.yaml 2>/dev/null || warn "请手动在 pubspec.yaml 的 flutter 块中添加 assets 配置"
    fi
fi

success "依赖更新完成"

# ========== 获取依赖 ==========
info "获取依赖..."
flutter pub get

# ========== 完成 ==========
echo ""
echo "=========================================="
echo -e "  ${GREEN}Flutter 项目初始化完成！${NC}"
echo "=========================================="
echo ""
echo "  项目名称: $PROJECT_NAME"
echo "  组织名称: $ORG_NAME"
echo "  项目路径: $SCRIPT_DIR"
echo ""
echo "  下一步:"
echo "  1. 修改 lib/config/env.dart 中的 API 地址"
echo "  2. 修改 lib/app.dart 中的应用名称"
echo "  3. 运行: flutter run"
echo ""
echo "  目录结构:"
echo "  lib/"
echo "  ├── config/    配置（环境、主题、路由）"
echo "  ├── core/      核心（网络、存储、工具、组件）"
echo "  ├── features/  业务模块（auth、home、profile）"
echo "  └── l10n/      国际化"
echo ""
echo "=========================================="
