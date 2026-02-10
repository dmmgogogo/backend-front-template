# Controller 设计模式

## BaseController 模式

### 为什么需要 BaseController？
- 统一响应格式
- 统一错误处理
- 复用用户信息提取
- 复用权限检查

---

## 标准 BaseController 实现

```go
package controllers

import (
    "encoding/json"
    "github.com/beego/beego/v2/core/logs"
    "github.com/beego/beego/v2/server/web"
    "github.com/beego/i18n"
)

type BaseController struct {
    web.Controller
    Lang       string      // 语言
    UserInfo   *User       // 当前用户
    DeviceID   string      // 设备ID
    Token      string      // JWT Token
    Code       int         // 响应码
    Msg        string      // 响应消息
    Result     interface{} // 响应数据
}

// Prepare 所有方法执行前调用
func (c *BaseController) Prepare() {
    // 提取语言
    c.Lang = c.Ctx.Input.Header("Accept-Language")
    if c.Lang == "" {
        c.Lang = "zh-CN"
    }
    i18n.SetLocale(c.Lang)
    
    // 提取设备ID
    c.DeviceID = c.Ctx.Input.Header("Device-Id")
    
    // 提取 Token（中间件已验证）
    c.Token = c.Ctx.Input.Header("token")
    
    // 从 Context 获取用户信息（JWT中间件写入）
    if userID := c.Ctx.Input.GetData("user_id"); userID != nil {
        c.UserInfo = &User{ID: userID.(int64)}
        // 从数据库加载完整用户信息
        c.UserInfo.LoadFromDB()
    }
}

// Success 成功响应
func (c *BaseController) Success(obj interface{}) {
    c.Code = 200
    c.Msg = "success"
    c.Result = obj
    c.TraceJson()
}

// Error 错误响应
func (c *BaseController) Error(code int, extraMsg ...string) {
    c.Code = code
    
    // 从 i18n 文件读取错误信息
    c.Msg = c.Tr("error." + fmt.Sprint(code))
    
    // 追加额外信息
    if len(extraMsg) > 0 && extraMsg[0] != "" {
        c.Msg = c.Msg + "(" + extraMsg[0] + ")"
    }
    
    c.Result = nil
    c.TraceJson()
}

// TraceJson 输出 JSON 并停止执行
func (c *BaseController) TraceJson() {
    c.Data["json"] = map[string]interface{}{
        "code": c.Code,
        "msg":  c.Msg,
        "data": c.Result,
    }
    c.ServeJSON()
    c.StopRun()
}

// ParseJson 解析 JSON 请求体
func (c *BaseController) ParseJson(v interface{}) error {
    err := json.Unmarshal(c.Ctx.Input.RequestBody, v)
    if err != nil {
        logs.Error("[ParseJson]Failed to parse: %v", err)
    }
    return err
}

// GetCurrentUserID 获取当前用户ID
func (c *BaseController) GetCurrentUserID() int64 {
    if c.UserInfo != nil {
        return c.UserInfo.ID
    }
    return 0
}

// Tr i18n 翻译
func (c *BaseController) Tr(key string, args ...interface{}) string {
    return i18n.Tr(c.Lang, key, args...)
}
```

---

## 业务 Controller 继承

### ✅ 正确做法
```go
package controllers

type UserController struct {
    BaseController  // 继承 BaseController
}

func (c *UserController) Login() {
    var req dto.LoginReq
    
    if err := c.ParseJson(&req); err != nil {
        c.Error(conf.PARAMS_ERROR, "参数解析失败")
        return
    }
    
    // 业务逻辑...
    
    c.Success(map[string]interface{}{
        "token": token,
        "user_info": userInfo,
    })
}
```

### ❌ 错误做法
```go
// 不要直接继承 web.Controller
type UserController struct {
    web.Controller  // ❌ 错误
}

// 不要手动构造响应
func (c *UserController) Login() {
    c.Data["json"] = map[string]interface{}{  // ❌ 错误
        "Code": 401,
        "Msg":  "未登录",
    }
    c.ServeJSON()
}
```

---

## 管理员 BaseController

### 扩展场景
当管理端需要特殊逻辑时，创建 `admin.BaseController`：

```go
package admin

import "your-project/controllers"

type BaseController struct {
    controllers.BaseController  // 继承通用 BaseController
}

// Prepare 管理员特定预处理
func (c *BaseController) Prepare() {
    // 调用父类 Prepare
    c.BaseController.Prepare()
    
    // 检查是否是管理员
    if !c.UserInfo.IsAdmin {
        c.Error(conf.FORBIDDEN, "需要管理员权限")
        return
    }
    
    // 记录管理员操作日志
    c.LogAdminAction()
}

func (c *BaseController) LogAdminAction() {
    // 记录操作日志到数据库
}
```

### 管理员 Controller 使用
```go
package admin

type AdminUserController struct {
    BaseController  // 继承管理员 BaseController
}

func (c *AdminUserController) List() {
    // 自动检查管理员权限
    // 自动记录操作日志
    
    users := GetAllUsers()
    c.Success(users)
}
```

---

## 权限控制模式

### 方法级权限检查
```go
func (c *BaseController) RequirePermission(permCode string) {
    hasPermission := CheckUserPermission(c.GetCurrentUserID(), permCode)
    if !hasPermission {
        c.Error(conf.FORBIDDEN, "无权限")
        return
    }
}

// 使用
func (c *EmployeeController) Delete() {
    c.RequirePermission("employee:delete")  // 检查权限
    
    // 业务逻辑...
    c.Success(nil)
}
```

### 装饰器模式（扩展）
```go
// 权限装饰器
func WithPermission(permCode string) func(*BaseController) {
    return func(c *BaseController) {
        c.RequirePermission(permCode)
    }
}
```

---

## 多租户模式

### 租户隔离 BaseController
```go
type BaseController struct {
    web.Controller
    MerchantID int64  // 租户ID
}

func (c *BaseController) Prepare() {
    // 从 JWT 提取租户ID
    merchantID := c.Ctx.Input.GetData("merchant_id")
    c.MerchantID = merchantID.(int64)
}

// 获取租户过滤的 ORM 查询
func (c *BaseController) GetTenantQuery(tableName string) orm.QuerySeter {
    db := orm.NewOrm()
    return db.QueryTable(tableName).Filter("merchant_id", c.MerchantID)
}
```

### 使用
```go
func (c *EmployeeController) List() {
    var employees []*Employee
    
    // 自动按租户过滤
    _, err := c.GetTenantQuery("employees").All(&employees)
    
    c.Success(employees)
}
```

---

## Finish 方法（清理）

### 资源清理
```go
func (c *BaseController) Finish() {
    // 记录请求日志
    logs.Info("[%s] %s %s - %d - %dms",
        c.Ctx.Input.Method(),
        c.Ctx.Input.URL(),
        c.GetCurrentUserID(),
        c.Code,
        time.Since(c.StartTime).Milliseconds(),
    )
    
    // 关闭数据库连接（如果有）
    // 清理临时文件
}
```

---

## 最佳实践

### ✅ 推荐做法
1. **所有 Controller 继承 BaseController**
2. **统一使用 `c.Success()` 和 `c.Error()` 返回响应**
3. **在 `Prepare()` 中处理通用逻辑**（用户信息、语言、设备ID）
4. **权限检查集中在 BaseController**
5. **多租户数据隔离在 BaseController 实现**

### ❌ 避免做法
1. 不要直接继承 `web.Controller`
2. 不要手动构造 JSON 响应
3. 不要在每个方法中重复提取用户信息
4. 不要在业务 Controller 中写权限检查逻辑

---

## 错误码管理

### 集中定义错误码
```go
// conf/error_code.go
package conf

const (
    SUCCESS      = 200
    PARAMS_ERROR = 400
    UNAUTHORIZED = 401
    FORBIDDEN    = 403
    NOT_FOUND    = 404
    SERVER_ERROR = 500
    
    // 业务错误码
    USER_NOT_EXIST   = 2001
    PASSWORD_ERROR   = 2002
    TOKEN_EXPIRED    = 2003
)
```

### i18n 错误消息
```ini
# conf/locale_zh.ini
[error]
200 = 成功
400 = 参数错误
401 = 未授权
403 = 禁止访问
500 = 服务器错误
2001 = 用户不存在
2002 = 密码错误
2003 = Token已过期

# conf/locale_en.ini
[error]
200 = Success
400 = Invalid parameters
401 = Unauthorized
403 = Forbidden
500 = Server error
2001 = User not found
2002 = Incorrect password
2003 = Token expired
```

### 使用
```go
c.Error(conf.UNAUTHORIZED)                    // "未授权"
c.Error(conf.PARAMS_ERROR, "邮箱格式不正确")    // "参数错误(邮箱格式不正确)"
```

---

## 参考实现
- 完整代码示例：`controllers/base.go`
- 管理员控制器：`controllers/admin/base.go`

