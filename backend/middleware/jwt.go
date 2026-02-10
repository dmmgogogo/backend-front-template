package middleware

import (
	"e-woms/conf"
	"e-woms/models"
	"slices"
	"std-library-slim/json"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

func JWTMiddleware(ctx *context.Context) {
	// 排除登录注册接口
	path := ctx.Request.URL.Path
	logs.Debug("JWTMiddleware path: %s", path)
	if slices.Contains(conf.NonLoginPathsBackend, path) || slices.Contains(conf.NonLoginPathsAdmin, path) {
		return
	}

	// 获取token (支持 token header 和 Authorization: Bearer header)
	authHeader := ctx.Input.Header("token")
	if authHeader == "" {
		// 尝试从 Authorization header 获取
		authHeader = ctx.Input.Header("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// 检查token是否在黑名单中
	if authHeader == "" || models.IsTokenBlacklisted(authHeader) {
		handleUnauthorized(ctx, "token无效")
		return
	}

	// 解析token
	claims, err := models.ParseJWTToken(authHeader)
	if err != nil {
		logs.Error("ParseJWTToken error: %v", err)
		handleUnauthorized(ctx, "无效的token")
		return
	}

	logs.Debug("解析token里面的内容: authHeader: %v, claims: %v", authHeader, json.String(claims))

	// 将用户信息存储在context中
	// JWT claims中的数字默认是float64,需要转换为int64
	if userID, ok := claims["user_id"]; ok {
		if userIDFloat, ok := userID.(float64); ok {
			ctx.Input.SetData("user_id", int64(userIDFloat))
		}
	}
	ctx.Input.SetData("username", claims["username"])

	if is_admin, ok := claims["is_admin"]; ok {
		// JWT claims中的数字默认是float64,需要转换为int
		if isAdminFloat, ok := is_admin.(float64); ok {
			ctx.Input.SetData("is_admin", int(isAdminFloat))
		}
	}
}

func handleUnauthorized(ctx *context.Context, msg string) {
	ctx.Output.SetStatus(401)
	ctx.Output.JSON(map[string]interface{}{
		"code": conf.UNAUTHORIZED,
		"msg":  msg,
	}, false, false)
}
