# Backend-Front Template

全栈项目通用模板，包含后端 + 前端（App & 管理后台）三个独立可选模块。

## 模块说明

| 模块 | 路径 | 技术栈 | 说明 |
|------|------|--------|------|
| 后端 | `backend/` | Go / Beego v2 / MySQL / Redis | 用户端 API + 管理端 API |
| App 前端 | `frontend/flutter/` | Flutter / Dart | 移动端应用 |
| 管理后台 | `frontend/vben-admin/` | Vue3 / Vben Admin / TypeScript | 管理后台前端 |

> 各模块独立，新项目按需选取，不需要的直接删除即可。

## 快速开始

### 1. 克隆模板

```bash
git clone <repo-url> my-new-project
cd my-new-project
```

### 2. 初始化后端

```bash
cd backend
cp conf/app.example.conf conf/app.conf
# 修改 app.conf 中的 MySQL / Redis 配置
# 修改 go.mod 中的 module 名称
# 修改 conf/const.go 中的 salt 和密钥
go run main.go
```

### 3. 初始化前端（按需）

**Flutter App：**

```bash
cd frontend/flutter
bash init.sh my_app com.yourcompany
```

**Vben Admin 管理后台：**

```bash
cd frontend/vben-admin
bash init.sh
```

## 后端 API 结构

| 命名空间 | 路径 | 用途 |
|---------|------|------|
| 管理端 | `/api/admin/*` | 管理后台接口（JWT + IP 白名单 + RBAC） |
| 用户端 | `/api/backend/*` | App 用户接口（JWT） |
| 通用 | `/api/common/*` | 文件上传等（JWT） |

统一响应格式：

```json
{ "code": 200, "msg": "success", "data": {} }
```

## AI 开发规范

每个模块携带自包含的 AI 规范文件（`CLAUDE.md`），复制模块时规范跟着走：

```
backend/CLAUDE.md              # 后端开发规范
frontend/flutter/CLAUDE.md     # Flutter 开发规范
frontend/vben-admin/CLAUDE.md  # Vben Admin 开发规范
.claude/                       # 通用 AI 技能库 & 指南
```

## 目录结构

```
├── .claude/                   # AI 规范 & 技能库
│   ├── README.md              # Agent vs Skills 指南
│   ├── SKILL.md               # Skills 索引
│   └── skills/                # 通用技能（beego-api / mysql / redis）
│
├── backend/                   # Go 后端
│   ├── CLAUDE.md
│   ├── controllers/           # 控制器（admin / backend / common）
│   ├── dto/                   # 数据传输对象
│   ├── models/                # 数据模型
│   ├── middleware/            # 中间件（CORS / JWT / 权限）
│   ├── services/              # 业务服务
│   └── conf/                  # 配置
│
├── frontend/
│   ├── flutter/               # Flutter App
│   │   ├── CLAUDE.md
│   │   └── init.sh            # 一键初始化
│   └── vben-admin/            # 管理后台
│       ├── CLAUDE.md
│       └── init.sh            # 一键初始化
│
└── README.md
```
