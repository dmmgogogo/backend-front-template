# Beego Models 目录规范

本目录包含 Beego 项目中 models 层的开发规范和最佳实践。

## 规范文件列表

- `01-directory-structure.md` - **Models 目录结构规范** ⭐ 必读

## 核心规则速查

### 📁 目录结构
```
models/
├── backend/          # API 后台模型（必须）
├── admin/            # 管理员后台模型（必须）
├── func.go           # 通用工具函数（允许）
└── token.go          # 通用结构（允许）
```

### ✅ 允许
- 业务模型放在 `models/api/` 或 `models/admin/`
- 工具函数放在 `models/func.go`
- 通用结构放在 `models/token.go`

### ❌ 禁止
- 业务模型直接放在 `models/` 根目录
- 在根目录创建业务实体文件如 `video.go`、`user.go`

### 📦 导入方式
```go
// API 后台
import apiModel "video-app-ai/models/backend"

// 管理员后台
import "video-app-ai/models/admin"
```

## 使用指南

1. **创建新模型前，先判断用途**
   - 服务于 API 接口 → `models/backend/`
   - 服务于管理后台 → `models/admin/`
   - 通用工具函数 → `models/` 根目录

2. **阅读详细规范**
   - 查看 `01-directory-structure.md` 了解完整规则和示例

3. **保持一致性**
   - 遵循现有项目的导入方式
   - 参考 `controllers/backend/base.go` 和 `controllers/admin/base.go` 的导入示例

## 注意事项

⚠️ **这是项目级别的强制规范，所有开发者必须遵守**

- 这个规范确保了代码结构清晰、职责分明
- 违反此规范会导致代码混乱、难以维护
- 新加入的开发者请务必先阅读此规范

## 快速决策树

```
创建新的模型文件
    ↓
是业务实体模型吗？
    ├─ 是 → 主要给谁用？
    │       ├─ API 用户端 → models/backend/
    │       └─ 管理员后台 → models/admin/
    └─ 否 → 是工具函数/通用结构吗？
            ├─ 是 → models/ 根目录
            └─ 否 → 考虑是否属于其他层（dto/services）
```

## 相关规范

- 参考 `beego-api/` - API 控制器规范
- 参考 `dto-usage.md` - DTO 层使用规范

---

**最后更新**: 2026-01-08
**维护者**: 项目团队

