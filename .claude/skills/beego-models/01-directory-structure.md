# Beego 项目 Models 目录结构规范

## 核心规范

**在 Beego 项目中，models/ 目录必须按照后台类型划分子目录，禁止直接在 models/ 根目录下创建业务模型文件。**

## 目录结构

```
models/
├── backend/              # API 后台模型目录
│   ├── user.go
│   ├── video.go
│   └── ...
├── admin/            # 管理员后台模型目录
│   ├── user.go
│   ├── admin_log.go
│   └── ...
├── func.go           # 通用工具函数（可放根目录）
└── token.go          # 通用token结构（可放根目录）
```

## 规则说明

### 1. 强制分目录

- **models/backend/** - 存放 API 后台相关的所有业务模型
- **models/admin/** - 存放管理员后台相关的所有业务模型

### 2. 允许根目录的文件类型

以下类型的文件可以放在 models/ 根目录：
- **通用工具函数**: 如 `func.go`（包含各后台共用的工具方法）
- **通用数据结构**: 如 `token.go`（JWT token 相关结构）
- **常量定义**: 如 `const.go`（全局常量）

### 3. 禁止根目录的文件类型

以下类型的文件**不允许**直接放在 models/ 根目录：
- ❌ 业务实体模型（如 User、Video、Order 等）
- ❌ 数据库表结构体
- ❌ 业务逻辑相关的数据模型
- ❌ API 请求/响应结构体

### 4. 导入规范

#### API 后台导入
```go
package backend

import (
    apiModel "video-app-ai/models/backend"
)

type BaseController struct {
    UserInfo *apiModel.User
}
```

#### 管理员后台导入
```go
package admin

import (
    "video-app-ai/models/admin"
)

type BaseController struct {
    UserInfo *admin.User
}
```

## 创建模型时的判断标准

### 问题1: 这个模型应该放在哪个目录？

**判断方法**:
1. 如果模型主要服务于 **API 接口**（用户端、移动端、第三方） → `models/backend/`
2. 如果模型主要服务于 **管理后台**（管理员、运营人员） → `models/admin/`
3. 如果是**两个后台都用的工具函数/常量** → `models/` 根目录（但要谨慎）

### 问题2: 两个后台都需要同一个实体怎么办？

**推荐方案**:
1. **分别创建** - 在 `models/backend/` 和 `models/admin/` 各创建一个版本（推荐）
   - 优点：解耦，可以根据不同后台定制字段
   - 例如：`models/backend/user.go` 和 `models/admin/user.go`

2. **提取到公共包** - 如果确实完全一致，可考虑创建 `models/common/`
   ```
   models/
   ├── common/          # 真正共用的模型（谨慎使用）
   │   └── video.go
   ├── backend/
   └── admin/
   ```

## 实际案例

### ✅ 正确示例

**场景**: 创建 AI 聊天功能的会话模型

**正确做法**:
```bash
# 判断：AI聊天是给用户（API端）使用的
# 创建在 models/backend/ 目录
models/backend/chat_session.go
```

```go
package backend

type ChatSession struct {
    SessionID string
    UserID    int64
    Messages  []ChatMessage
}
```

**在控制器中使用**:
```go
package backend

import (
    apiModel "video-app-ai/models/backend"
)

func (c *AIChatController) Chat() {
    session := &apiModel.ChatSession{}
}
```

### ❌ 错误示例

**场景**: 创建xx模型

**错误做法**:
```bash
# ❌ 直接放在根目录
models/xx.go
```

**正确做法**:
```bash
# ✅ 根据用途放到对应目录
models/backend/xx.go      # 如果是API端展示用
# 或
models/admin/xx.go    # 如果是管理后台用
# 或两个都创建
```

## 命名规范

### 文件命名
- 使用下划线分隔: `chat_session.go`, `user_profile.go`
- 文件名应描述模型用途，清晰易懂

### 包名
- `models/backend/` 目录下的文件，package 声明为 `package backend`
- `models/admin/` 目录下的文件，package 声明为 `package admin`

### 导入别名
- 导入 `models/backend` 时推荐使用别名: `backendModel "video-app-ai/models/backend"`


## 特殊说明

### DTO (Data Transfer Object) 的位置

DTO 层应该有独立的目录，不放在 models/ 下：
```
dto/
├── backend/              # API 后台的 DTO
│   ├── chat_request.go
│   └── chat_response.go
└── admin/            # 管理后台的 DTO
    ├── login_request.go
    └── user_response.go
```

### Service 层模型

如果某些数据结构只在 Service 层使用，不涉及数据库：
- 可以直接定义在 Service 文件中
- 或创建 `services/types/` 目录存放

## 总结检查清单

创建模型前，问自己：

- [ ] 这是业务模型还是工具函数？
- [ ] 这个模型主要服务于哪个后台？
- [ ] 是否应该放在 `models/backend/` 或 `models/admin/`？
- [ ] 包名和导入路径是否正确？
- [ ] 是否遵循了命名规范？

**记住：永远不要把业务模型直接放在 models/ 根目录！**

