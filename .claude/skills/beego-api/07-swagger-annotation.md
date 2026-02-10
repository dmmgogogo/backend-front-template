# Swagger 文档规范

## 为什么需要 Swagger？
- 自动生成 API 文档
- 前后端接口约定
- 在线测试 API
- 减少沟通成本

---

## 安装 Swagger 工具

### bee 工具（Beego官方）
```bash
# 安装 bee
go install github.com/beego/bee/v2@latest

# 生成 Swagger 文档
bee generate docs

# 生成位置：swagger/swagger.json
```

---

## 标准注解格式

### 完整注解模板
```go
// MethodName 方法简要描述
// @Title 接口标题
// @Description 接口详细描述
// @Param param_name param_type data_type required "参数说明"
// @Success 200 {object} response_type "成功响应示例"
// @Failure 400 {object} map[string]interface{} "错误描述"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/path [method]
func (c *Controller) MethodName() {
    // ...
}
```

---

## 参数类型说明

### @Param 语法
```
@Param 参数名 参数位置 数据类型 是否必填 "参数说明" [其他属性]
```

### 参数位置

| 参数位置 | 说明 | 示例 |
|---------|------|------|
| `body` | 请求体（JSON） | `@Param body body dto.LoginReq true "请求参数"` |
| `query` | 查询参数（GET） | `@Param page query int false "页码" default(1)` |
| `path` | 路径参数 | `@Param id path int true "用户ID"` |
| `header` | 请求头 | `@Param token header string true "认证Token"` |
| `formData` | 表单数据 | `@Param file formData file true "上传文件"` |

---

## 常见接口类型示例

### 1. 登录接口（POST + JSON Body）
```go
// Login 用户登录
// @Title 用户登录
// @Description 用户通过邮箱和密码进行登录
// @Param body body dto.LoginReq true "请求参数"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{"token":"xxx","user_info":{}}}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "邮箱或密码错误"
// @Failure 403 {object} map[string]interface{} "账号已被禁用"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/user/login [post]
func (c *UserController) Login() {
    var req dto.LoginReq
    // ...
}
```

### 2. 列表查询（GET + Query 参数）
```go
// List 订单列表
// @Title 订单列表查询
// @Description 查询用户订单列表，支持分页和筛选
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query int false "订单状态：1-待支付 2-已支付 3-已取消"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{"list":[],"total":0,"page":1,"page_size":20}}"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/order/list [get]
func (c *OrderController) List() {
    page, _ := c.GetInt("page", 1)
    // ...
}
```

### 3. 创建资源（POST + JSON Body）
```go
// Create 创建订单
// @Title 创建订单
// @Description 用户购买算力套餐，创建订单
// @Param body body dto.CreateOrderReq true "请求参数"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{"order_id":123,"order_no":"xxx"}}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/order/create [post]
func (c *OrderController) Create() {
    var req dto.CreateOrderReq
    // ...
}
```

### 4. 更新资源（PUT + Path 参数）
```go
// Update 更新用户资料
// @Title 更新用户资料
// @Description 更新用户个人资料信息
// @Param id path int true "用户ID"
// @Param body body dto.UpdateProfileReq true "请求参数"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{}}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/user/update/:id [put]
func (c *UserController) Update() {
    userId := c.Ctx.Input.Param(":id")
    // ...
}
```

### 5. 删除资源（DELETE + Path 参数）
```go
// Delete 删除订单
// @Title 删除订单
// @Description 删除指定订单（仅未支付订单可删除）
// @Param id path int true "订单ID"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{"id":123}}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 404 {object} map[string]interface{} "订单不存在"
// @Failure 403 {object} map[string]interface{} "订单已支付，无法删除"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/order/delete/:id [delete]
func (c *OrderController) Delete() {
    orderId := c.Ctx.Input.Param(":id")
    // ...
}
```

### 6. 文件上传（POST + formData）
```go
// Upload 文件上传
// @Title 文件上传
// @Description 上传文件到服务器（支持图片、文档）
// @Param file formData file true "上传的文件"
// @Param type formData string false "文件类型：image-图片 document-文档" default(image)
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{"url":"xxx","filename":"xxx","size":1024}}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 413 {object} map[string]interface{} "文件过大"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/common/upload [post]
func (c *CommonController) Upload() {
    file, header, err := c.GetFile("file")
    // ...
}
```

### 7. 需要 Token 认证的接口
```go
// Info 获取用户信息
// @Title 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Param token header string true "认证Token"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{"id":1,"email":"xxx","nickname":"xxx"}}"
// @Failure 401 {object} map[string]interface{} "Token无效或已过期"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/user/info [get]
func (c *UserController) Info() {
    userId := c.GetCurrentUserID()
    // ...
}
```

---

## 响应示例格式

### 成功响应
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [...],
    "total": 100
  }
}
```

### 错误响应
```json
{
  "code": 400,
  "msg": "参数错误(邮箱格式不正确)",
  "data": null
}
```

### Swagger 注解中表示
```go
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{...}}"
// @Failure 400 {object} map[string]interface{} "{"code":400,"msg":"参数错误","data":null}"
```

---

## Controller 注释（全局配置）

### main.go 注释
```go
// @APIVersion 1.0.0
// @Title Nexora 云算力投资平台 API
// @Description RESTful API 文档
// @Contact support@nexora.com
// @TermsOfServiceUrl http://nexora.com/terms/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package main
```

### router.go 分组注释
```go
// @Title 用户相关接口
// @Description 用户注册、登录、个人信息管理
```

---

## 生成 Swagger 文档

### 方法1：bee 工具
```bash
# 在项目根目录执行
bee generate docs

# 生成文件
swagger/
├── swagger.json    # Swagger 2.0 格式
└── swagger.yml     # YAML 格式
```

### 方法2：手动脚本
创建 `test/generate_swagger.sh`：
```bash
#!/bin/bash
# Swagger 文档生成脚本

echo "正在生成 Swagger 文档..."

# 安装 bee 工具（如果未安装）
if ! command -v bee &> /dev/null; then
    echo "安装 bee 工具..."
    go install github.com/beego/bee/v2@latest
fi

# 生成文档
bee generate docs

echo "✅ Swagger 文档生成完成: swagger/swagger.json"
```

```bash
chmod +x test/generate_swagger.sh
./test/generate_swagger.sh
```

---

## 查看 Swagger 文档

### 方法1：Swagger UI（推荐）
```bash
# 下载 Swagger UI
git clone https://github.com/swagger-api/swagger-ui.git
cd swagger-ui/dist

# 复制 swagger.json 到 dist 目录
cp /path/to/your/project/swagger/swagger.json ./

# 启动本地服务器
python3 -m http.server 8080

# 浏览器访问
# http://localhost:8080
# 输入 URL: http://localhost:8080/swagger.json
```

### 方法2：在线工具
1. 访问 https://editor.swagger.io/
2. 上传 `swagger/swagger.json`
3. 在线查看和测试

---

## 常见问题

### 问题1：DTO 不识别
**原因**：DTO 包未被引用

**解决**：在 `routers/router.go` 中引入
```go
import (
    _ "your-project/dto"  // 引入 DTO
)
```

### 问题2：注释不生效
**原因**：注释格式不正确

**检查**：
- `@router` 必须在最后一行
- HTTP 方法必须小写：`[get]`, `[post]`, `[put]`, `[delete]`
- 路径必须以 `/` 开头

### 问题3：生成的文档不完整
**解决**：
```bash
# 清除缓存重新生成
rm -rf swagger/
bee generate docs
```

---

## 最佳实践

### ✅ 推荐做法
1. **所有 API 方法必须有完整注释**
2. **使用 DTO 定义请求/响应结构**
3. **@Failure 覆盖所有可能的错误码**
4. **@Description 清晰描述业务场景**
5. **定期生成文档并同步给前端**

### ❌ 避免做法
1. ❌ 不写注释
2. ❌ 注释与实际代码不一致
3. ❌ 使用内联结构体而不是 DTO
4. ❌ 遗漏必填参数标记

---

## 注解速查表

| 注解 | 说明 | 示例 |
|------|------|------|
| `@Title` | 接口标题 | `@Title 用户登录` |
| `@Description` | 接口描述 | `@Description 用户通过邮箱登录` |
| `@Param` | 参数定义 | `@Param body body dto.LoginReq true "请求参数"` |
| `@Success` | 成功响应 | `@Success 200 {object} dto.LoginResp` |
| `@Failure` | 失败响应 | `@Failure 401 {object} map[string]interface{} "未授权"` |
| `@router` | 路由定义 | `@router /api/login [post]` |

---

## 参考资料
- Beego Swagger 文档: https://beego.wiki/docs/advantage/docs/
- Swagger 规范: https://swagger.io/specification/
- Swagger UI: https://github.com/swagger-api/swagger-ui

