# DTO（Data Transfer Object）使用规范

## 为什么需要 DTO？

### ❌ 不使用 DTO 的问题
```go
func (c *UserController) Login() {
    // 内联结构体 - 不可复用
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    json.Unmarshal(c.Ctx.Input.RequestBody, &req)
    
    // 无法在 Swagger 中复用
    // 无法在测试中复用
    // 修改时需要改多处
}
```

### ✅ 使用 DTO 的优势
1. **可复用** - 一次定义，多处使用
2. **可维护** - 统一修改，不会遗漏
3. **可文档化** - Swagger 自动识别
4. **可测试** - 测试代码可直接引用

---

## DTO 目录结构

```
dto/
├── user.go       # 用户相关请求（登录、注册等）
├── order.go      # 订单相关请求
├── recharge.go   # 充值相关请求
├── admin.go      # 管理员相关请求
├── package.go    # 套餐相关请求
├── level.go      # 等级相关请求
└── common.go     # 通用请求（分页等）
```

---

## DTO 命名规范

| 类型 | 命名规则 | 示例 |
|------|---------|------|
| 请求 DTO | 以 `Req` 结尾 | `LoginReq`, `RegisterReq` |
| 响应 DTO | 以 `Resp` 或 `Res` 结尾 | `LoginResp`, `UserInfoResp` |
| 查询参数 | 以 `Query` 结尾 | `ListQuery`, `SearchQuery` |

---

## DTO 定义示例

### user.go
```go
package dto

// SendCodeReq 发送验证码请求
type SendCodeReq struct {
    Email string `json:"email" valid:"Email,Required"`
}

// RegisterReq 用户注册请求
type RegisterReq struct {
    Email      string `json:"email" valid:"Email,Required"`
    Password   string `json:"password" valid:"MinSize(6);MaxSize(20)"`
    Code       string `json:"code" valid:"Required"`
    InviteCode string `json:"invite_code"` // 可选
}

// LoginReq 用户登录请求
type LoginReq struct {
    Email    string `json:"email" valid:"Email,Required"`
    Password string `json:"password" valid:"Required"`
}

// UpdateProfileReq 更新个人资料请求
type UpdateProfileReq struct {
    Nickname string `json:"nickname"`
    Avatar   string `json:"avatar"`
    Phone    string `json:"phone" valid:"Mobile"`
}
```

### common.go
```go
package dto

// PageQuery 分页查询参数
type PageQuery struct {
    Page     int    `json:"page" form:"page" valid:"Min(1)"`
    PageSize int    `json:"page_size" form:"page_size" valid:"Range(1,100)"`
    Keyword  string `json:"keyword" form:"keyword"`
}

// IDReq 通用ID请求
type IDReq struct {
    ID int64 `json:"id" valid:"Required;Min(1)"`
}

// IDsReq 批量ID请求
type IDsReq struct {
    IDs []int64 `json:"ids" valid:"Required"`
}
```

### order.go
```go
package dto

// CreateOrderReq 创建订单请求
type CreateOrderReq struct {
    PackageID int64  `json:"package_id" valid:"Required;Min(1)"`
    Quantity  int    `json:"quantity" valid:"Min(1);Max(100)"`
    CouponID  *int64 `json:"coupon_id"` // 可选，使用指针表示可空
}

// OrderListQuery 订单列表查询
type OrderListQuery struct {
    PageQuery               // 嵌入分页参数
    Status    int    `json:"status" form:"status"`
    StartDate string `json:"start_date" form:"start_date"`
    EndDate   string `json:"end_date" form:"end_date"`
}
```

---

## Controller 中使用 DTO

### ✅ 正确用法
```go
package controllers

import "your-project/dto"

func (c *UserController) Login() {
    var req dto.LoginReq
    
    if err := c.ParseJson(&req); err != nil {
        c.Error(conf.PARAMS_ERROR, "参数解析失败")
        return
    }
    
    // 参数验证
    if err := c.Validate(req); err != nil {
        c.Error(conf.PARAMS_ERROR, err.Error())
        return
    }
    
    // 业务逻辑
    token, err := models.UserLogin(req.Email, req.Password)
    if err != nil {
        c.Error(conf.PASSWORD_ERROR)
        return
    }
    
    c.Success(map[string]interface{}{
        "token": token,
    })
}
```

### GET 请求使用 DTO
```go
func (c *OrderController) List() {
    var query dto.OrderListQuery
    
    // 解析 GET 参数
    if err := c.ParseForm(&query); err != nil {
        c.Error(conf.PARAMS_ERROR)
        return
    }
    
    // 默认值
    if query.Page == 0 {
        query.Page = 1
    }
    if query.PageSize == 0 {
        query.PageSize = 20
    }
    
    list, total := models.GetOrderList(c.GetCurrentUserID(), query)
    
    c.Success(map[string]interface{}{
        "list":  list,
        "total": total,
        "page":  query.Page,
    })
}
```

---

## 参数验证

### Beego 内置验证
```go
import "github.com/beego/beego/v2/core/validation"

// 在 BaseController 中添加验证方法
func (c *BaseController) Validate(obj interface{}) error {
    valid := validation.Validation{}
    ok, err := valid.Valid(obj)
    
    if err != nil {
        return err
    }
    
    if !ok {
        for _, err := range valid.Errors {
            return fmt.Errorf("%s: %s", err.Key, err.Message)
        }
    }
    
    return nil
}
```

### 验证标签
```go
type RegisterReq struct {
    Email    string `valid:"Email,Required"`           // 必填邮箱
    Password string `valid:"MinSize(6);MaxSize(20)"`   // 6-20字符
    Phone    string `valid:"Mobile"`                   // 手机号
    Age      int    `valid:"Range(18,100)"`            // 18-100
    URL      string `valid:"URL"`                      // URL格式
}
```

### 常用验证规则
| 规则 | 说明 | 示例 |
|------|------|------|
| `Required` | 必填 | `valid:"Required"` |
| `Email` | 邮箱格式 | `valid:"Email"` |
| `Mobile` | 手机号 | `valid:"Mobile"` |
| `MinSize(n)` | 最小长度 | `valid:"MinSize(6)"` |
| `MaxSize(n)` | 最大长度 | `valid:"MaxSize(20)"` |
| `Range(min,max)` | 数值范围 | `valid:"Range(1,100)"` |
| `Min(n)` | 最小值 | `valid:"Min(1)"` |
| `Max(n)` | 最大值 | `valid:"Max(999)"` |
| `URL` | URL格式 | `valid:"URL"` |

---

## 响应 DTO（可选）

### 定义响应结构
```go
package dto

// LoginResp 登录响应
type LoginResp struct {
    Token    string       `json:"token"`
    UserInfo UserInfoResp `json:"user_info"`
}

// UserInfoResp 用户信息响应
type UserInfoResp struct {
    ID       int64  `json:"id"`
    Email    string `json:"email"`
    Nickname string `json:"nickname"`
    Avatar   string `json:"avatar"`
    Level    int    `json:"level"`
}
```

### 使用响应 DTO
```go
func (c *UserController) Login() {
    var req dto.LoginReq
    c.ParseJson(&req)
    
    // 业务逻辑
    user, token := models.UserLogin(req.Email, req.Password)
    
    // 构造响应
    resp := dto.LoginResp{
        Token: token,
        UserInfo: dto.UserInfoResp{
            ID:       user.ID,
            Email:    user.Email,
            Nickname: user.Nickname,
            Avatar:   user.Avatar,
            Level:    user.Level,
        },
    }
    
    c.Success(resp)
}
```

---

## 嵌入式 DTO（复用）

### 嵌入分页参数
```go
type UserListQuery struct {
    dto.PageQuery              // 嵌入分页
    Status int    `form:"status"`
    Role   string `form:"role"`
}
```

### 嵌入时间范围
```go
type DateRangeQuery struct {
    StartDate string `json:"start_date" form:"start_date"`
    EndDate   string `json:"end_date" form:"end_date"`
}

type OrderListQuery struct {
    dto.PageQuery
    dto.DateRangeQuery  // 嵌入时间范围
    Status int `form:"status"`
}
```

---

## DTO 转换 Model

### 手动转换
```go
func (req *CreateOrderReq) ToModel(userID int64) *models.Order {
    return &models.Order{
        UserID:    userID,
        PackageID: req.PackageID,
        Quantity:  req.Quantity,
        Status:    1,
        CreatedAt: time.Now().Unix(),
    }
}

// 使用
order := req.ToModel(c.GetCurrentUserID())
db.Insert(order)
```

### 使用库转换（可选）
```bash
go get github.com/jinzhu/copier
```

```go
import "github.com/jinzhu/copier"

var order models.Order
copier.Copy(&order, &req)
order.UserID = c.GetCurrentUserID()
```

---

## Swagger 集成

### DTO 在 Swagger 中使用
```go
// Login 用户登录
// @Title 用户登录
// @Description 用户通过邮箱和密码进行登录
// @Param body body dto.LoginReq true "请求参数"
// @Success 200 {object} dto.LoginResp
// @router /api/user/login [post]
func (c *UserController) Login() {
    // ...
}
```

Swagger 会自动识别 `dto.LoginReq` 和 `dto.LoginResp` 结构，生成文档。

---

## 最佳实践

### ✅ 推荐做法
1. **所有请求参数定义在 `/dto` 目录**
2. **DTO 命名遵循规范**（Req/Resp/Query 结尾）
3. **使用验证标签**，避免手动验证
4. **复用嵌入式 DTO**（如 PageQuery）
5. **为复杂响应定义 DTO**

### ❌ 避免做法
1. ❌ 不要使用内联结构体
2. ❌ 不要在多处重复定义相同结构
3. ❌ 不要把 Model 直接当 DTO 用
4. ❌ 不要在 DTO 中写业务逻辑

---

## 示例项目结构
```
dto/
├── common.go      # PageQuery, IDReq
├── user.go        # LoginReq, RegisterReq, UpdateProfileReq
├── order.go       # CreateOrderReq, OrderListQuery
├── recharge.go    # RechargeReq
└── admin.go       # AdminLoginReq, CreateMerchantReq
```

---

## 参考资料
- Beego 验证文档: https://beego.wiki/docs/mvc/controller/validation/
- copier 库: https://github.com/jinzhu/copier

