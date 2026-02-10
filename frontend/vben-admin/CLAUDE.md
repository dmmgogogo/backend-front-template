# Vben Admin 管理后台 AI 开发规范

本文档为 Vben Admin 模块的自包含 AI 开发规范，复制到新项目时跟着走。

---

## 技术栈

| 项 | 选型 |
|----|------|
| 框架 | Vue 3.x |
| 构建 | Vite |
| UI 框架 | Ant Design Vue / Element Plus（取决于 Vben 版本） |
| 状态管理 | Pinia |
| 路由 | Vue Router 4 |
| 网络请求 | Axios |
| 国际化 | vue-i18n |
| 包管理 | pnpm（Monorepo） |
| 语言 | TypeScript |

---

## 项目结构

运行 `init.sh` 后生成的标准结构：

```
vben-admin/
├── apps/
│   └── web-admin/                     # 主应用
│       ├── src/
│       │   ├── api/                   # API 接口定义
│       │   │   ├── core/             # 请求封装
│       │   │   │   └── request.ts    # Axios 实例（对接后端 response 格式）
│       │   │   ├── system/           # 系统管理 API（/api/admin/*）
│       │   │   │   ├── user.ts       # 管理员接口
│       │   │   │   ├── role.ts       # 角色接口
│       │   │   │   └── menu.ts       # 菜单/权限接口
│       │   │   └── business/         # 业务 API（按模块划分）
│       │   │
│       │   ├── views/                # 页面视图
│       │   │   ├── _core/            # 核心页（登录、404）
│       │   │   ├── dashboard/        # 仪表盘
│       │   │   ├── system/           # 系统管理
│       │   │   │   ├── user/         # 用户管理
│       │   │   │   ├── role/         # 角色管理
│       │   │   │   └── menu/         # 菜单管理
│       │   │   └── business/         # 业务页（按模块划分）
│       │   │
│       │   ├── router/               # 路由配置
│       │   │   ├── routes/           # 路由表（模块化）
│       │   │   └── guard/            # 路由守卫
│       │   │
│       │   ├── store/                # 状态管理（Pinia）
│       │   │   └── modules/
│       │   │       ├── user.ts
│       │   │       └── permission.ts
│       │   │
│       │   ├── locales/              # 国际化
│       │   │   └── langs/
│       │   │       ├── zh-CN.ts
│       │   │       └── en.ts
│       │   │
│       │   └── utils/                # 工具函数
│       │
│       ├── .env                      # 通用环境变量
│       ├── .env.development          # 开发环境
│       ├── .env.production           # 生产环境
│       └── package.json
│
├── packages/                         # 共享包（Monorepo）
├── pnpm-workspace.yaml
└── package.json
```

---

## 初始化

```bash
# 使用初始化脚本（推荐）
bash init.sh

# 脚本会自动：检查环境 → 克隆官方仓库 → 配置 .env → 创建业务目录 → 安装依赖
```

---

## 与后端对接规范

### API 基础地址

```
管理端 API: /api/admin/*
通用 API:   /api/common/*
```

> 用户端 API (`/api/backend/*`) 由 Flutter App 对接，管理后台不调用。

### 统一响应格式

后端返回：`{ "code": 200, "msg": "success", "data": {...} }`

```typescript
interface BackendResponse<T = any> {
  code: number;
  msg: string;
  data: T;
}
```

### Axios 封装要点

```
1. baseURL → 从 .env 读取 VITE_GLOB_API_URL + VITE_GLOB_API_URL_PREFIX
2. 请求拦截器 → 自动注入 token header
3. 响应拦截器:
   - code === 200 → 直接返回 data
   - code !== 200 → 提示 msg，特定 code 做特殊处理
   - 401 → 清除登录态 → 跳转登录页
4. 错误处理 → 网络异常、超时统一 toast 提示
```

### Token 管理

```
请求头（与后端 middleware/jwt.go 对齐）:
  token: <jwt_token>
  或 Authorization: Bearer <jwt_token>

存储: localStorage 或 cookie
退出: 调用 POST /api/admin/user/logout（后端将 token 加入黑名单）
```

### 文件上传

```
POST /api/common/upload
Content-Type: multipart/form-data
字段名: file
需携带 token
```

---

## 环境变量

```bash
# .env.development
VITE_GLOB_API_URL=http://localhost:8282
VITE_GLOB_API_URL_PREFIX=/api/admin

# .env.production
VITE_GLOB_API_URL=https://api.yourcompany.com
VITE_GLOB_API_URL_PREFIX=/api/admin
```

---

## 权限对接规范

### 登录流程

```
POST /api/admin/user/login
  → 返回 token + 基础用户信息
  → 存储 token

GET /api/admin/user/userinfo
  → 返回完整用户信息 + 角色列表 + 权限列表
  → 前端根据权限动态生成路由和菜单
```

### 权限控制层级

| 层级 | 实现方式 |
|------|---------|
| 路由级 | 根据后端权限列表过滤路由表 |
| 菜单级 | 动态生成侧边栏菜单 |
| 按钮级 | `v-auth` 指令控制显隐 |
| API 级 | 后端 `middleware/permission.go` 兜底 |

### 路由 & 菜单约定

```typescript
/**
 * 后端 permission path  →  前端路由 path
 * /api/admin/user/list  →  /system/user
 * /api/admin/role/list  →  /system/role
 * /api/admin/xxx/list   →  /business/xxx
 *
 * 路由 meta:
 * - title: 菜单名称（i18n key）
 * - icon: 菜单图标
 * - permission: 后端权限标识
 * - hideMenu: 是否隐藏
 * - orderNo: 排序
 */
```

---

## 新模块开发流程

```bash
# 以「订单管理」模块为例

# 1. 创建 API 文件
# apps/web-admin/src/api/business/order.ts

# 2. 创建页面
# apps/web-admin/src/views/business/order/index.vue
# apps/web-admin/src/views/business/order/detail.vue

# 3. 添加路由
# apps/web-admin/src/router/routes/modules/order.ts

# 4. 添加国际化
# apps/web-admin/src/locales/langs/zh-CN.ts
# apps/web-admin/src/locales/langs/en.ts
```

### API 文件模板

```typescript
// api/business/order.ts
import { request } from '../core/request';

enum Api {
  List = '/order/list',
  Detail = '/order/detail',
  Create = '/order/create',
  Update = '/order/update',
  Delete = '/order/delete',
}

export function getOrderList(params: any) {
  return request.get({ url: Api.List, params });
}

export function getOrderDetail(id: number) {
  return request.get({ url: Api.Detail, params: { id } });
}
```

---

## 命名规范

```
文件名:    kebab-case       (order-list.vue)
组件名:    PascalCase       (OrderList)
变量/函数:  camelCase        (getOrderList)
类/接口:   PascalCase       (OrderModel)
常量:      UPPER_SNAKE_CASE (API_PREFIX)
CSS class: kebab-case       (order-container)
```
