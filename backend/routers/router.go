package routers

import (
	"e-woms/controllers/admin"
	"e-woms/controllers/backend"
	"e-woms/controllers/common"
	"e-woms/middleware"

	"github.com/beego/beego/v2/server/web"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouters() {
	web.InsertFilter("*", web.BeforeRouter, middleware.Cors)
	web.InsertFilter("/api/*", web.BeforeRouter, middleware.JWTMiddleware)

	// 添加首页
	web.Router("/", &backend.MainController{})

	// 管理平台
	ns := web.NewNamespace("/api",
		// 管理员系统
		web.NSNamespace("/admin",
			web.NSRouter("/user/login", &admin.UserController{}, "post:Login"),
			web.NSRouter("/user/logout", &admin.UserController{}, "post:Logout"),
			web.NSRouter("/user/userinfo", &admin.UserController{}, "get:GetUserInfo"),
			web.NSRouter("/user/change-password", &admin.UserController{}, "post:ChangePassword"),
		),

		// 公开接口
		web.NSNamespace("/backend",
			// 基础模块
			web.NSRouter("/user/send-code", &backend.UserController{}, "post:SendCode"),
			web.NSRouter("/user/forgot-password", &backend.UserController{}, "post:ForgotPassword"),
			web.NSRouter("/user/register", &backend.UserController{}, "post:Register"),
			web.NSRouter("/user/login", &backend.UserController{}, "post:Login"),
			web.NSRouter("/user/logout", &backend.UserController{}, "post:Logout"),
			web.NSRouter("/user/userinfo", &backend.UserController{}, "get:GetUserInfo"),
		),

		// 通用功能
		web.NSNamespace("/common",
			// 文件上传（需登录）
			web.NSRouter("/upload", &common.UploadController{}, "post:Upload"),
		),
	)
	web.AddNamespace(ns)

	// Swagger文档路由 (开发模式下启用)
	if web.BConfig.RunMode == "dev" {
		web.Handler("/swagger/*", httpSwagger.WrapHandler)
	}
}
