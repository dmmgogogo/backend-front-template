# Beego v2 框架基础

## ⚠️ 重要：Beego v2 vs v1

**本技能库基于 Beego v2，不是 Beego v1。千万不要搞混！**

## 版本对比

### 导入路径对比

#### ✅ Beego v2（正确）
```go
import "github.com/beego/beego/v2/server/web"
import "github.com/beego/beego/v2/client/orm"
import "github.com/beego/beego/v2/core/logs"
import "github.com/beego/beego/v2/server/web/context"
```

#### ❌ Beego v1（错误 - 不要使用）
```go
// 错误 - 已废弃
import "github.com/astaxie/beego"
import "github.com/astaxie/beego/orm"
import "github.com/astaxie/beego/logs"
```

### API 对比表

| 功能 | Beego v2 (正确) | Beego v1 (错误) |
|------|----------------|-----------------|
| Controller | `web.Controller` | `beego.Controller` |
| 配置 | `web.AppConfig.String()` | `beego.AppConfig.String()` |
| 路由 | `web.Router()` | `beego.Router()` |
| 运行 | `web.Run()` | `beego.Run()` |
| ORM | `orm.NewOrm()` | `orm.NewOrm()` (同名但包路径不同) |
| 日志 | `logs.Error()` | `logs.Error()` (同名但包路径不同) |

---

## 项目初始化

### go.mod 配置
```go
module your-project

go 1.18

require (
    github.com/beego/beego/v2 v2.1.0
    // 其他依赖...
)
```

### main.go 标准模板
```go
package main

import (
    "github.com/beego/beego/v2/core/logs"
    "github.com/beego/beego/v2/server/web"
    _ "your-project/routers"  // 引入路由
)

func main() {
    // 初始化日志/语言
    initLocales()
    initLogger()
    // 初始化路由
	routers.InitRouters()
	// 初始化mysql
	services.InitMysql()
	// 初始化redis
	services.InitRedis()
    // 初始化定时任务
    initCronshell()

    // 启动服务
    web.Run()
}

// 初始化多语言
func initLocales() {
	langs := []string{"zh", "en"}
	for _, lang := range langs {
		logs.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			logs.Error("Fail to set message file: " + err.Error())
			return
		}
	}
}

// 初始化log
func initLogger() {
	logs.SetLogger(logs.AdapterFile, `{ "filename": "logs/app.log", "daily": true, "maxdays": 7}`)

	// logs.SetLogger(logs.AdapterConsole)
	logs.SetLevel(logs.LevelDebug) // 设置日志级别为 Debug
	logs.SetLogFuncCall(true)      // 显示文件名和行号
	logs.SetLogFuncCallDepth(3)
}

// 初始化定时任务
func initCronshell() {
    // 创建 cron 实例，使用秒级精度
	cronInstance := cron.New(cron.WithSeconds())

	// Cron 表达式格式：秒 分 时 日 月 周
	_, err := cronInstance.AddFunc("0 */1 * * * *", func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("[Cron][initCronshell] panic recovered: %v", r)
			}
		}()

		logs.Info("[Cron][initCronshell] ========== START ==========")

		// TODO 执行任务

		logs.Info("[Cron][initCronshell] ========== COMPLETED ==========")
	})

	if err != nil {
		logs.Error("[Cron][InitCron] failed to register initCronshell task: %v", err)
		panic(err)
	}

	// 启动定时任务
	cronInstance.Start()
	logs.Info("[Cron][initCronshell] cron scheduler started successfully")
}
```

---

## 配置文件 (conf/app.conf)

### 基础配置
```ini
# 应用配置
appname = your-project
httpport = 8080
runmode = dev

# 日志配置
logs_level = debug

# 数据库配置（JSON格式）
MYSQL_CONFIG = {"host":"127.0.0.1:3306","user":"root","password":"xxx","db":"mydb"}
REDIS_CONFIG = {"addr":"127.0.0.1:6379","password":"","db":0}

# JWT 配置
JWT_SECRET = your-secret-key
JWT_SECRET_EXPIRE_TIME = 8760  # 365天（小时）

# 限流配置
RATE_LIMIT_SECOND_CREATE_NUM = 100
RATE_LIMIT_COMMON_NUM = 500
```

### 读取配置
```go
import "github.com/beego/beego/v2/server/web"

// 读取字符串
appName := web.AppConfig.String("appname")

// 读取整数（带默认值）
httpPort, _ := web.AppConfig.Int("httpport")

// 读取 JSON 配置
mysqlConfig := web.AppConfig.DefaultString("MYSQL_CONFIG", "{}")
```

---

## 路由注册 (routers/router.go)

### 命名空间路由（推荐）
```go
package routers

import (
    "github.com/beego/beego/v2/server/web"
    "your-project/controllers"
)

func init() {
    // API 命名空间
    apiNS := web.NewNamespace("/api",
        // 用户相关
        web.NSRouter("/user/login", &controllers.UserController{}, "post:Login"),
        web.NSRouter("/user/register", &controllers.UserController{}, "post:Register"),
        web.NSRouter("/user/info", &controllers.UserController{}, "get:Info"),
        
        // 订单相关
        web.NSRouter("/order/list", &controllers.OrderController{}, "get:List"),
        web.NSRouter("/order/create", &controllers.OrderController{}, "post:Create"),
    )
    
    web.AddNamespace(apiNS)
}
```

### 简单路由
```go
func init() {
    web.Router("/", &controllers.MainController{})
    web.Router("/api/login", &controllers.UserController{}, "post:Login")
}
```

---

## Controller 基础

### 标准 Controller 结构
```go
package controllers

import (
    "github.com/beego/beego/v2/server/web"
)

type UserController struct {
    web.Controller
}

// Prepare 在方法执行前调用
func (c *UserController) Prepare() {
    // 获取请求头
    token := c.Ctx.Input.Header("token")
    
    // 预处理逻辑
}

// Login 登录方法
func (c *UserController) Login() {
    // 获取参数
    email := c.GetString("email")
    password := c.GetString("password")
    
    // 业务逻辑...
    
    // 返回 JSON
    c.Data["json"] = map[string]interface{}{
        "code": 200,
        "msg":  "success",
        "data": map[string]interface{}{
            "token": "xxx",
        },
    }
    c.ServeJSON()
}
```

---

## 中间件注册

### 全局中间件
```go
import "github.com/beego/beego/v2/server/web"

func init() {
    // CORS 中间件（所有路由）
    web.InsertFilter("*", web.BeforeRouter, CorsMiddleware)
    
    // JWT 认证中间件（排除登录/注册）
    web.InsertFilter("/api/*", web.BeforeRouter, JWTMiddleware, 
        web.WithExcludePatterns("/api/login", "/api/register"))
}
```

### 中间件实现
```go
import "github.com/beego/beego/v2/server/web/context"

func CorsMiddleware(ctx *context.Context) {
    ctx.Output.Header("Access-Control-Allow-Origin", "*")
    ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
    
    if ctx.Input.Method() == "OPTIONS" {
        ctx.Output.SetStatus(200)
    }
}
```

---

## 常用方法

### 获取请求参数
```go
// GET/POST 参数
email := c.GetString("email")
page, _ := c.GetInt("page", 1)         // 默认值 1
status, _ := c.GetBool("status", true)

// 路径参数 /user/:id
userId := c.Ctx.Input.Param(":id")

// 请求头
token := c.Ctx.Input.Header("token")

// JSON Body 解析
var req struct {
    Email string `json:"email"`
}
json.Unmarshal(c.Ctx.Input.RequestBody, &req)
```

### 返回响应
```go
// JSON 响应
c.Success({
    "code": 200,
    "msg":  "success",
    "data": data,
})
```

---

## ORM 注册

### 模型注册
```go
package models

import "github.com/beego/beego/v2/client/orm"

type User struct {
    ID       int64  `orm:"pk;auto"`
    Email    string `orm:"size(100)"`
    Password string `orm:"size(255)"`
}

func init() {
    // 注册模型
    orm.RegisterModel(new(User))
}

func (u *User) TableName() string {
    return "app_users"
}
```

---

## 最佳实践

### ✅ 推荐做法
1. **统一导入路径**：始终使用 `github.com/beego/beego/v2/*`
2. **配置外部化**：敏感信息放在 `app.conf` 或环境变量
3. **命名空间路由**：大型项目使用 `web.NewNamespace`
4. **中间件复用**：认证、日志、限流等逻辑抽取为中间件

### ❌ 避免做法
1. 不要混用 v1 和 v2 包
2. 不要在 Controller 中直接写 SQL
3. 不要硬编码敏感信息（密码、密钥）
4. 不要在 `main.go` 中写业务逻辑

---

## 参考资料
- [Beego v2 官方文档](https://beego.wiki/)
- [迁移指南 v1 → v2](https://beego.wiki/docs/intro/upgrade/)

