## 通用技能库（可复用到其他项目）

### Beego API 开发
- `.claude/skills/beego-api/00-README.md` - 导航索引
- `.claude/skills/beego-api/01-framework-basics.md` - Beego v2 框架基础
- `.claude/skills/beego-api/02-controller-pattern.md` - Controller 设计模式
- `.claude/skills/beego-api/05-dto-usage.md` - DTO 使用规范
- `.claude/skills/beego-api/07-swagger-annotation.md` - Swagger 文档规范

### Beego Models
- `.claude/skills/beego-models/00-README.md` - 导航索引
- `.claude/skills/beego-models/01-directory-structure.md` - Model 目录结构

### MySQL 数据库
- `.claude/skills/mysql/00-README.md` - 导航索引
- `.claude/skills/mysql/01-orm-basic.md` - Beego ORM v2 基础操作

### Redis 缓存
- `.claude/skills/redis/00-README.md` - 导航索引
- `.claude/skills/redis/01-basic-operations.md` - Redis 基础操作

---

## 各模块独立 AI 规范

每个模块的 CLAUDE.md 是自包含的，取用模块时规范跟着走：

| 模块 | AI 规范 | 初始化脚本 |
|------|---------|-----------|
| 后端 Go/Beego v2 | `backend/CLAUDE.md` | - |
| Flutter App | `frontend/flutter/CLAUDE.md` | `frontend/flutter/init.sh` |
| Vben Admin | `frontend/vben-admin/CLAUDE.md` | `frontend/vben-admin/init.sh` |

---

## 项目特定内容（根据不同项目定制）

- `.claude/agent/context/business-modules.md` - 业务模块说明
- `.claude/agent/context/code-standards.md` - 项目编码规范
- `CLAUDE.md` - 项目约定（自动加载）

## 依赖本地其他组件库
- std-library-slim

---
