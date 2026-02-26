# AGENTS.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Overview

全栈项目通用模板，包含三个独立可选模块：
- `backend/` — Go 1.24 + Beego v2 后端（用户端 API + 管理端 API）
- `frontend/flutter/` — Flutter 3.x 移动端 App
- `frontend/vben-admin/` — Vue 3 + Vben Admin 管理后台（pnpm monorepo）

各模块完全独立，新项目按需选取，不需要的直接删除。

## Commands

### Backend (Go / Beego v2)

```bash
cd backend

# Run
go run main.go

# Run tests (uses goconvey)
go test ./tests/...

# Run a single test
go test ./tests/ -run TestBeego -v

# Build
go build -o app main.go
```

Config: copy `conf/app.example.conf` → `conf/app.conf` and update MySQL/Redis credentials before first run.

Backend runs on port **8282** in dev mode. Swagger UI available at `/swagger/` in dev mode.

### Flutter App

```bash
cd frontend/flutter

# First-time init (creates Flutter project in place)
bash init.sh <project_name> [org_name]
# e.g.: bash init.sh my_app com.yourcompany

# After init, from the generated Flutter project root:
flutter pub get
flutter run

# Run tests
flutter test

# Build
flutter build apk       # Android
flutter build ios       # iOS
```

### Vben Admin (Vue 3 / pnpm monorepo)

```bash
cd frontend/vben-admin

# First-time init (clones official Vben repo and configures it)
bash init.sh

# After init, from the generated vben-admin directory:
pnpm install
pnpm dev        # runs apps/web-admin
pnpm build      # production build
pnpm lint       # lint
pnpm typecheck  # TypeScript check
```

## Architecture

### Backend Request Flow

```
HTTP Request
  → middleware/cors.go       (all routes)
  → middleware/jwt.go        (/api/* — sets user_id in context)
  → controller Prepare()     (auth checks, i18n setup)
  → controller method
  → c.Success(data) / c.Error(code)
  → JSON: { "code": 200, "msg": "...", "data": {} }
```

**Two distinct BaseControllers** with different auth strategies:
- `controllers/backend/base.go` — parses JWT from `Authorization: Bearer <token>`, sets `c.UserId`
- `controllers/admin/base.go` — reads `user_id` injected by JWT middleware, enforces IP whitelist, loads full `*admin.User` into `c.UserInfo`, provides `LogOperation()` for audit trail

Controllers call `c.Error(code)` which halts execution (`StopRun()`), so no explicit returns are needed after error calls.

### API Namespaces

| Namespace | Prefix | Auth |
|-----------|--------|------|
| 管理端 | `/api/admin/*` | JWT + IP whitelist + RBAC (`middleware/permission.go`) |
| 用户端 | `/api/backend/*` | JWT only |
| 通用 | `/api/common/*` | JWT (file upload) |

Non-login paths are whitelisted in `conf/permission.go` (`NonLoginPathsAdmin`, `NonLoginPathsBackend`).

### Backend Module: Adding a Feature

1. **Model** — `models/backend/<name>.go` or `models/admin/<name>.go`. Table prefix `app_`. Register in `init()` via `orm.RegisterModel(new(Model))`. Primary key `id bigint`, timestamps `created_time`/`updated_time` as Unix bigint (`orm:"auto_now_add;type(bigint)"`).
2. **DTO** — `dto/backend/<name>.go` or `dto/admin/<name>.go`. Use `valid:` struct tags for validation.
3. **Controller** — embed the appropriate `BaseController`. Use `c.ParseJson(&req)`, `c.Success(data)`, `c.Error(conf.SOME_CODE)`.
4. **Route** — register in `routers/router.go` under the correct namespace.

### go.mod Module Name

The module is currently named `e-woms`. When using this as a template for a new project, update `go.mod` (`module <newname>`) and all import paths. Also note:

```
replace std-library-slim => ../../../GO/src/go-company/std-library-slim/
```

This `replace` directive points to a local path outside the repo. The admin `base.go` uses `std-library-slim/json`. Adjust or remove this dependency when adapting the template.

### Flutter App Architecture

After `init.sh`, the app follows a feature-first structure under `lib/features/<feature>/{data,models,providers,pages}`. Core infrastructure lives in `lib/core/` (network, storage, utils, widgets).

- HTTP client is Dio-based (`lib/core/network/http_client.dart`) with auth and error interceptors.
- State management: Provider (simple) or Riverpod (complex).
- Logging: all log calls go through `lib/core/log/app_logger.dart` — never use `print`/`debugPrint` directly in business code.
- Any operation >300ms must show visible loading feedback with interaction disabled (button `onPressed: null`).

Flutter only calls `/api/backend/*` and `/api/common/*`. Admin API is Vben Admin's domain.

### Vben Admin Architecture

After `init.sh`, the main app lives at `apps/web-admin/`. API calls target `/api/admin/*`. Axios is configured via `.env.development` / `.env.production` with `VITE_GLOB_API_URL` + `VITE_GLOB_API_URL_PREFIX`.

Permission model: backend returns a permission list at login; frontend filters routes and menu dynamically. Button-level control uses `v-auth` directive. Backend RBAC in `middleware/permission.go` is the final guard.

### Auth Token Convention (all clients)

```
Header (either):
  token: <jwt_token>
  Authorization: Bearer <jwt_token>

Logout → POST /api/[admin|backend]/user/logout  (adds token to Redis blacklist)
401 response → clear local token → redirect to login
```

Language negotiation: backend reads `Accept-Language` header (backend API) or `Language` header (admin API). Supports `zh` and `en`.

## AI Context Files

Each module carries its own detailed AI dev spec:
- `backend/CLAUDE.md` — backend conventions, error codes, iOS IAP integration
- `frontend/flutter/CLAUDE.md` — Flutter patterns, logging rules, IAP client side
- `frontend/vben-admin/CLAUDE.md` — Vben Admin structure, permission wiring

Reusable technical skills (framework usage, ORM patterns, Redis ops) are in `.claude/skills/`. Project-specific business context goes in `.claude/agent/context/`.

## Git Branching

```
main      → production (protected)
develop   → main dev branch
feature/* → from develop
hotfix/*  → from main
release/* → from develop
```
