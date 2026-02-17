# Flutter App AI 开发规范

本文档为 Flutter 模块的自包含 AI 开发规范，复制到新项目时跟着走。

---

## 技术栈

| 项 | 选型 |
|----|------|
| 框架 | Flutter 3.x |
| 语言 | Dart 3.x |
| 网络请求 | Dio |
| 状态管理 | Provider（轻量）/ Riverpod（复杂） |
| 路由 | go_router |
| 本地存储 | SharedPreferences |
| JSON 序列化 | json_serializable + build_runner |
| 国际化 | flutter_localizations + intl |

---

## 项目结构

运行 `init.sh` 后生成的标准结构：

```
lib/
├── main.dart                      # 入口文件
├── app.dart                       # App 根组件（MaterialApp / 路由 / 主题）
│
├── config/                        # 全局配置
│   ├── env.dart                   # 环境变量（baseUrl、appName）
│   ├── theme.dart                 # 主题配色
│   └── routes.dart                # 路由表定义
│
├── core/                          # 核心层（不含业务逻辑）
│   ├── network/                   # 网络请求
│   │   ├── http_client.dart       # Dio 封装
│   │   ├── api_response.dart      # 统一响应模型
│   │   └── interceptors/          # 拦截器
│   │       ├── auth_interceptor.dart    # Token 注入
│   │       └── error_interceptor.dart   # 错误处理
│   ├── storage/                   # 本地存储
│   │   └── sp_util.dart           # SharedPreferences 封装
│   ├── utils/                     # 工具类
│   └── widgets/                   # 通用基础组件
│
├── features/                      # 业务模块（按功能拆分）
│   ├── auth/                      # 认证模块
│   │   ├── data/                  # API 接口 + Repository
│   │   ├── models/                # 数据模型
│   │   ├── providers/             # 状态管理
│   │   └── pages/                 # 页面
│   ├── home/
│   └── profile/
│
├── l10n/                          # 国际化资源
│   ├── intl_en.arb
│   └── intl_zh.arb
│
└── generated/                     # 自动生成文件

assets/                            # 静态资源
├── images/
├── icons/
└── fonts/
```

---

## 初始化

```bash
# 使用初始化脚本（推荐）
bash init.sh <project_name> [org_name]

# 示例
bash init.sh my_app com.yourcompany
```

脚本会自动：创建 Flutter 项目 → 生成目录结构 → 创建基础文件 → 配置 pubspec.yaml

---

## 基础依赖

```yaml
dependencies:
  flutter:
    sdk: flutter
  dio: ^5.x                  # 网络请求
  provider: ^6.x             # 状态管理
  go_router: ^14.x           # 路由管理
  shared_preferences: ^2.x   # 本地存储
  json_annotation: ^4.x      # JSON 注解
  logger: ^2.x               # 日志
  flutter_localizations:      # 国际化
    sdk: flutter
  intl: any

dev_dependencies:
  json_serializable: ^6.x    # JSON 代码生成
  build_runner: ^2.x         # 构建工具
  flutter_test:
    sdk: flutter
```

---

## 与后端对接规范

### API 基础地址

```
用户端 API: /api/backend/*
通用 API:   /api/common/*
```

> 管理后台 API (`/api/admin/*`) 由 Vben Admin 对接，Flutter 不调用。

### 统一响应模型

后端返回格式：`{ "code": 200, "msg": "success", "data": {...} }`

```dart
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
```

### Dio 封装原则

```
1. baseUrl 从 EnvConfig 读取
2. 请求拦截器 → 自动添加 token header
3. 响应拦截器 → code != 200 时统一处理
4. 错误拦截器 → 网络异常、超时统一提示
5. 401 → 清除本地 token → 跳转登录页
```

### Token 管理

```
请求头（与后端 middleware/jwt.go 对齐）:
  token: <jwt_token>
  或 Authorization: Bearer <jwt_token>

流程:
  登录成功 → 持久化存储 token
  每次请求 → 拦截器自动携带 token
  401 响应 → 清除 token → 跳登录
  主动退出 → 调用 POST /api/backend/user/logout（后端将 token 加入黑名单）
```

### 国际化

```
请求头: Accept-Language: zh | en
后端根据此 header 返回对应语言的错误信息
前端切换语言时同步更新请求 header
```

### 文件上传

```
POST /api/common/upload
Content-Type: multipart/form-data
字段名: file
需携带 token
```

---

## 环境配置

```dart
enum Env { dev, staging, prod }

class EnvConfig {
  static Env current = Env.dev;

  static String get baseUrl {
    switch (current) {
      case Env.dev:     return 'http://localhost:8282/api/backend';
      case Env.staging: return 'https://staging-api.yourcompany.com/api/backend';
      case Env.prod:    return 'https://api.yourcompany.com/api/backend';
    }
  }

  static bool get isDebug => current == Env.dev;
}
```

---

## 新模块开发流程

```bash
# 1. 创建模块目录
mkdir -p lib/features/order/{data,models,providers,pages}

# 2. 创建文件（按此顺序）
# models/order_model.dart   → 数据模型
# data/order_api.dart       → API 接口定义
# data/order_repository.dart → 数据仓库
# providers/order_provider.dart → 状态管理
# pages/order_list_page.dart    → 页面

# 3. 注册路由（config/routes.dart）
# 4. 添加国际化文本（l10n/intl_*.arb）
```

---

## 命名规范

```
文件名:   snake_case        (user_model.dart)
变量/函数: camelCase         (getUserInfo)
类名:     PascalCase        (UserModel)
常量:     lowerCamelCase    (defaultTimeout)
私有成员:  _camelCase        (_isLoading)
```

---

## iOS 支付（内购验单）

本模板通过 **init.sh** 可选注入 **iOS App Store 内购** 代码，用于赞助/打赏等场景：

- **位置**: 若存在 `patches/ios_support_service.dart`，init 会复制为 `lib/features/support/support_service.dart`，并添加 `in_app_purchase` 依赖。
- **后端接口**: `POST /api/backend/support/ios/verify`（需登录），见 **backend CLAUDE.md** 的「iOS 支付」一节。
- **用法**: `await SupportService().purchaseIOSSupport(amount: 10, onDebugEvent: (msg) => debugPrint(msg));`
- **商品 ID**: `support_service.dart` 内 `_iosProductAmount` 需与后端 `controllers/backend/support_ios.go` 的 `iosProductAmount` 一致；路径为 `/support/ios/verify`，请求由 `HttpClient` 的 Auth 拦截器自动带 token。

**相关文件**（便于 AI 定位“iOS 支付代码”）：

- `patches/ios_support_service.dart` — 模板源码，init 时复制到 `lib/features/support/support_service.dart`
- `lib/features/support/support_service.dart` — 运行时的 iOS 内购 + 验单逻辑（init 后存在）
