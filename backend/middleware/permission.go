package middleware

import (
	"e-woms/conf"
	"e-woms/models/admin"
	"slices"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

// 请求 /api/enterprise/role/list
// ↓
// 1. 检查是否是登录接口 → 否
// ↓
// 2. 获取用户信息 (user_id, merchant_id, is_admin)
// ↓
// 3. 检查是否是管理员 → 如果是,直接通过
// ↓
// 4. 查询用户的所有权限
// ↓
// 5. 匹配路由和HTTP方法
// ↓
// 6. 有权限 → 通过 | 无权限 → 返回403

// PermissionMiddleware 权限验证中间件
func PermissionMiddleware(ctx *context.Context) {
	// 0. 排除登录接口，不需要权限验证
	if slices.Contains(conf.NoPermissionCheckPathsBackend, ctx.Request.URL.Path) {
		return
	}

	// 1. 从Token获取用户信息
	userIDData := ctx.Input.GetData("user_id")
	usernameData := ctx.Input.GetData("username")
	isAdminData := ctx.Input.GetData("is_admin")

	userID, userIDOk := userIDData.(int64)
	username, usernameOk := usernameData.(string)
	isAdmin, isAdminOk := isAdminData.(int)

	logs.Debug("[PermissionMiddleware] Type assertions: userIDOk=%v, usernameOk=%v, isAdminOk=%v", userIDOk, usernameOk, isAdminOk)
	logs.Debug("[PermissionMiddleware] userID=%v, username=%s, isAdmin=%d", userID, username, isAdmin)

	// 如果缺少必要的用户信息,返回未登录错误
	if !userIDOk || !usernameOk || !isAdminOk {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{
			"code": 401,
			"msg":  "未登录或登录信息无效",
		}, false, false)
		return
	}

	// 2. 企业管理员跳过权限检查
	// if isAdmin == 1 {
	// 	logs.Debug("[PermissionMiddleware] 企业管理员跳过权限检查: userID=%d, username=%s, isAdmin=%d", userID, username, isAdmin)
	// 	return
	// }

	// 3. 获取当前请求的路由和方法
	route := ctx.Request.URL.Path
	method := ctx.Request.Method

	// 4. 查询用户是否有该路由权限
	hasPermission := checkUserPermission(userID, route, method)

	// 5. 无权限返回403
	if !hasPermission {
		ctx.Output.SetStatus(403)
		ctx.Output.JSON(map[string]interface{}{
			"code": 403,
			"msg":  "无权限访问",
		}, false, false)
		return
	}
}

// checkUserPermission 检查用户是否有指定路由的权限
func checkUserPermission(userID int64, route string, method string) bool {
	// 查询用户的所有权限ID
	rolePermissionModel := &admin.RolePermission{}
	permissionIDs, err := rolePermissionModel.GetUserPermissions(userID)
	if err != nil || len(permissionIDs) == 0 {
		return false
	}

	// 批量查询权限详情
	permissionModel := &admin.Permission{}
	permissions, err := permissionModel.GetByIDs(permissionIDs)
	if err != nil {
		return false
	}

	// 检查是否有匹配的权限
	for _, perm := range permissions {
		// 简单的路由匹配:检查权限的API路由是否与当前请求路由匹配
		// 注意:这里使用简单的字符串匹配,实际项目中可能需要更复杂的路由匹配逻辑
		if matchRoute(perm.APIRoute, route) && perm.HTTPMethod == method {
			return true
		}
	}

	return false
}

// matchRoute 路由匹配函数
// 支持简单的路径参数匹配,例如 /api/role/detail/:id 匹配 /api/role/detail/1
func matchRoute(pattern string, path string) bool {
	// 简单实现:如果pattern中包含:参数,则进行前缀匹配
	// 例如: /api/enterprise/role/detail/:id 可以匹配 /api/enterprise/role/detail/1

	// 去除末尾的斜杠
	if len(pattern) > 0 && pattern[len(pattern)-1] == '/' {
		pattern = pattern[:len(pattern)-1]
	}
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// 完全匹配
	if pattern == path {
		return true
	}

	// 检查是否包含路径参数
	patternParts := splitPath(pattern)
	pathParts := splitPath(path)

	// 长度必须相同
	if len(patternParts) != len(pathParts) {
		return false
	}

	// 逐段匹配
	for i := range patternParts {
		// 如果是路径参数(以:开头),则跳过
		if len(patternParts[i]) > 0 && patternParts[i][0] == ':' {
			continue
		}
		// 否则必须完全匹配
		if patternParts[i] != pathParts[i] {
			return false
		}
	}

	return true
}

// splitPath 分割路径
func splitPath(path string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if i > start {
				parts = append(parts, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		parts = append(parts, path[start:])
	}
	return parts
}
