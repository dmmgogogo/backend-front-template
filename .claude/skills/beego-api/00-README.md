# Beego v2 API 开发技能库

## 📚 目录导航

### 核心框架
- **[01-framework-basics.md](01-framework-basics.md)** - Beego v2 框架基础（v1 vs v2 区别）
- **[02-controller-pattern.md](02-controller-pattern.md)** - Controller 设计模式
- **[03-middleware.md](03-middleware.md)** - 中间件开发
- **[04-routing.md](04-routing.md)** - 路由配置

### 数据交互
- **[05-dto-usage.md](05-dto-usage.md)** - DTO 请求参数管理
- **[06-response-format.md](06-response-format.md)** - 统一响应格式

### 文档与测试
- **[07-swagger-annotation.md](07-swagger-annotation.md)** - Swagger 文档规范
- **[08-error-handling.md](08-error-handling.md)** - 错误处理最佳实践

### 认证与权限
- **[09-jwt-auth.md](09-jwt-auth.md)** - JWT 认证实现
- **[10-rbac-permission.md](10-rbac-permission.md)** - RBAC 权限模型

---

## 🎯 适用范围
此技能库适用于所有基于 **Beego v2** 的 RESTful API 项目

## 📦 技术栈要求
- Go 1.22+
- Beego v2 (`github.com/beego/beego/v2`)

## 🔄 使用方式

### 新项目初始化
```bash
# 复制技能库到新项目
cp -r .claude/skills/beego-api /path/to/new-project/.claude/skills/

# 复制 agent 模板
cp -r .claude/templates/agent-template /path/to/new-project/.claude/agent/
```

### 项目特定配置
在 `.claude/agent/` 中覆盖或扩展技能库内容

---

## 📝 更新日志
- **2026-01-06**: 初始创建，拆分自 SKILL.md

