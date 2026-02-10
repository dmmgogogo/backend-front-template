package conf

// 权限相关配置

// 非登录path - 管理平台
var NonLoginPathsAdmin = []string{
	"/api/admin/user/login",
	"/api/ip-manage", // IP白名单管理接口
}

// 非登录path - 前端平台
var NonLoginPathsBackend = []string{
	"/api/backend/user/send-code",
	"/api/backend/user/register",
	"/api/backend/user/login",
	"/api/common/upload",
}

// 不需要权限校验的接口 - 前端平台
var NoPermissionCheckPathsBackend = []string{
	"/api/backend/user/send-code",
	"/api/backend/user/register",
	"/api/backend/user/login",
	"/api/backend/user/forgot-password",
	"/api/backend/user/change-password",
	"/api/common/upload",
}
