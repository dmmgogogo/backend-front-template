# 后端 AI 开发规范（Go / Beego v2）

本文档为后端模块的自包含 AI 开发规范，复制到新项目时跟着走。

---

## 技术栈

| 项 | 选型 |
|----|------|
| 语言 | Go 1.24+ |
| 框架 | Beego v2.3+ |
| ORM | Beego ORM v2 |
| 数据库 | MySQL |
| 缓存 | Redis |
| 认证 | JWT (golang-jwt/jwt/v4) |
| 文档 | Swagger (swaggo/http-swagger) |
| 国际化 | beego/i18n |

---

## 项目结构

```
backend/
├── main.go                    # 入口：初始化 → 路由 → MySQL → Redis → 启动
├── conf/                      # 配置
│   ├── app.conf               # 运行时配置（不入 Git）
│   ├── app.example.conf       # 配置模板（入 Git）
│   ├── const.go               # 常量（token 密钥、salt 等）
│   ├── error_code.go          # 统一错误码定义
│   ├── permission.go          # 权限路径配置
│   ├── locale_zh.ini          # 中文翻译
│   └── locale_en.ini          # 英文翻译
│
├── controllers/               # HTTP 控制器
│   ├── admin/                 # 管理后台 API（/api/admin/*）
│   │   └── base.go            # 管理端 BaseController（含 IP 白名单、操作日志）
│   ├── backend/               # 用户端 API（/api/backend/*）
│   │   └── base.go            # 用户端 BaseController（含 JWT 解析、i18n）
│   └── common/                # 通用 API（/api/common/*）
│
├── dto/                       # 数据传输对象
│   ├── admin/                 # 管理端 DTO
│   └── backend/               # 用户端 DTO
│
├── models/                    # 数据模型 & ORM
│   ├── admin/                 # 管理端模型（app_admin_users、app_roles 等）
│   ├── backend/               # 用户端模型（app_users）
│   ├── func.go                # 通用模型函数
│   └── token.go               # JWT 令牌黑名单
│
├── middleware/                # 中间件
│   ├── cors.go                # 跨域处理
│   ├── jwt.go                 # JWT 认证
│   └── permission.go          # RBAC 权限校验
│
├── routers/                   # 路由
│   └── router.go              # 主路由配置
│
├── services/                  # 业务服务
│   ├── init.go                # MySQL & Redis 初始化
│   ├── mail.go                # 邮件服务
│   └── password_validator.go  # 密码校验
│
└── utils/                     # 工具函数
    ├── jwt_util.go            # JWT 工具
    ├── google_authenticator.go # 2FA（可选）
    └── ip_whitelist.go        # IP 白名单
```

---

## API 命名空间

| 命名空间 | 路径前缀 | 用途 | 认证 |
|---------|---------|------|------|
| 管理端 | `/api/admin/*` | 管理后台 API | JWT + IP 白名单 + RBAC |
| 用户端 | `/api/backend/*` | App/用户 API | JWT |
| 通用 | `/api/common/*` | 文件上传等 | JWT |

---

## 统一响应格式

所有 API 统一返回此 JSON 格式：

```json
{
  "code": 200,
  "msg": "success",
  "data": {}
}
```

### 错误码约定

| code | 含义 |
|------|------|
| 200 | 成功 |
| 401 | 未认证 / token 过期 |
| 403 | 无权限 |
| 500 | 服务器内部错误 |
| 其他 | 见 `conf/error_code.go` |

错误信息支持 i18n，根据请求头 `Accept-Language: zh|en` 返回对应语言。

---

## Token 规范

```
请求头携带方式（二选一）：
  token: <jwt_token>
  Authorization: Bearer <jwt_token>

管理端 token 包含 is_admin: true 标记
退出登录时 token 加入 Redis 黑名单
```

---

## 新模块开发流程

### 1. 添加新的 Controller

```go
// controllers/backend/order.go
package backend

type OrderController struct {
    BaseController  // 继承用户端基础控制器
}

func (c *OrderController) GetList() {
    // 业务逻辑
    c.Success(data)
}
```

### 2. 添加 DTO

```go
// dto/backend/order.go
package backend

type CreateOrderReq struct {
    ProductId int64  `json:"product_id" valid:"Required"`
    Quantity  int    `json:"quantity" valid:"Required;Min(1)"`
}
```

### 3. 添加 Model

```go
// models/backend/order.go
package backend

type Order struct {
    Id          int64  `orm:"auto;pk" json:"id"`
    UserId      int64  `orm:"index" json:"user_id"`
    ProductId   int64  `json:"product_id"`
    Status      int    `orm:"default(0)" json:"status"`
    CreatedTime int64  `orm:"auto_now_add;type(bigint)" json:"created_time"`
    UpdatedTime int64  `orm:"auto_now;type(bigint)" json:"updated_time"`
}

func init() {
    orm.RegisterModel(new(Order))
}

func (o *Order) TableName() string {
    return "app_orders"
}
```

### 4. 注册路由

```go
// routers/router.go — 在对应命名空间下添加
web.NSRouter("/order/list", &backend.OrderController{}, "get:GetList"),
```

---

## 数据库约定

- 表名前缀: `app_`
- 主键: `id` (auto increment bigint)
- 时间戳: `created_time` / `updated_time` (Unix timestamp bigint)
- 软删除: `deleted_time`（可选）
- Model 在 `init()` 中注册 ORM

---

## 配置管理

```bash
# 新项目初始化
cp conf/app.example.conf conf/app.conf
# 然后修改：
# 1. MySQL 连接信息
# 2. Redis 连接信息
# 3. conf/const.go 中的 salt 和 JWT 密钥
```

---

## iOS 支付（内购验单）

本模板已内置 **iOS App Store 内购验单** 能力，用于赞助/打赏等场景：

- **接口**: `POST /api/backend/support/ios/verify`（需登录，Header 带 `Authorization: Bearer <token>`）
- **请求体**: `product_id`、`transaction_id`、`receipt_data`（App Store 收据 base64）
- **流程**: 校验参数 → 向 Apple 验单（生产 + 沙盒 21007 回退）→ 按 `transaction_id` 幂等写入 `app_support_orders` → 更新用户 `support_total_amount`、`support_level`、`vip`

**相关文件**（便于 AI 定位“iOS 支付代码”）：

- `controllers/backend/support_ios.go` — 验单接口、`iosProductAmount` 商品映射、`verifyAppleReceipt` 调 Apple
- `models/backend/support_order.go` — 订单表模型
- `models/backend/support.go` — `AddSupportByTransaction`、`resolveSupportLevel`
- `dto/backend/support.go` — 请求/响应 DTO
- `docs/014_ios_support.sql` — 用户表字段 + `app_support_orders` 建表

**配置**: `app.conf` 或 `app.example.conf` 中 `ios_iap_shared_secret`（App Store Connect 内购“App 专用共享密钥”）。商品 ID 与金额映射在 `support_ios.go` 的 `iosProductAmount` 中，按项目修改。

---

## 相关 Skills

详细技术文档参见根目录 `.claude/skills/`：
- `beego-api/` — Beego v2 API 开发模式
- `beego-models/` — Model 目录结构
- `mysql/` — ORM 增删改查
- `redis/` — Redis 缓存操作
