package admin

import (
	"e-woms/conf"
	"e-woms/dto/admin"
	adminDto "e-woms/dto/admin"
	"e-woms/models"
	adminModel "e-woms/models/admin"
	"std-library-slim/json"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

// UserController 认证控制器
type UserController struct {
	BaseController
}

// Login 管理员登录
// @Summary 管理员登录
// @Description 通过邮箱和密码登录，返回JWT Token和管理员列表
// @Tags 后台-用户管理
// @Accept json
// @Produce json
// @Param body body LoginForm true "登录表单"
// @Success 200 {object} map[string]interface{} "{"code": 200, "msg": "success", "data": {"token": "jwt_token", "user": {...}, "default_merchant": {...}, "merchants": [...]}}"
// @Failure 401 {object} map[string]interface{} "{"code": 401, "msg": "邮箱或密码错误"}"
// @router /api/admin/user/login [post]
func (c *UserController) Login() {
	var form adminDto.LoginForm
	err := c.ParseJson(&form)
	if err != nil {
		logs.Error("[UserController][Login] unmarshal error: %v", err)
		c.Error(conf.PARAMS_ERROR, "参数解析失败")
		return
	}

	// 参数校验
	if form.Username == "" {
		c.Error(conf.PARAMS_ERROR, "邮箱或用户名不能为空")
		return
	}
	if form.Password == "" {
		c.Error(conf.PARAMS_ERROR, "密码不能为空")
		return
	}
	if form.VerifyCode == "" {
		c.Error(conf.PARAMS_ERROR, "Google验证码不能为空")
		return
	}

	logs.Debug("[UserController][Login] form: %v", json.String(form))

	// 用户登录验证（使用邮箱）
	adminInfo := &adminModel.User{}
	err = adminInfo.LoginByUsername(form.Username, form.Password)
	if err != nil {
		logs.Error("[UserController][Login] login error: %v", err)
		c.Error(conf.UNAUTHORIZED, "邮箱或用户名或密码错误")
		return
	}

	// // 校验 Google 验证码
	// if adminInfo.VerifyCode == "" {
	// 	logs.Error("[UserController][Login] user %s has no Google Authenticator secret", form.Username)
	// 	c.Error(conf.UNAUTHORIZED, "账户未绑定Google验证码")
	// 	return
	// }
	// if !utils.VerifyGoogleAuthCode(adminInfo.VerifyCode, form.VerifyCode) {
	// 	logs.Error("[UserController][Login] user %s Google Authenticator verification failed", form.Username)
	// 	c.Error(conf.UNAUTHORIZED, "Google验证码错误")
	// 	return
	// }

	logs.Debug("[UserController][Login] user: %v", json.String(adminInfo))

	// 生成 JWT Token
	token, err := models.GenerateAdminJWTToken(adminInfo.ID, adminInfo.Username)
	if err != nil {
		logs.Error("[AdminLogin]Failed to generate token: %v", err)
		c.Error(conf.SERVER_ERROR, "生成Token失败")
		return
	}

	// 查询用户的角色
	userRoleModel := &adminModel.UserRole{}
	roles, err := userRoleModel.GetUserRoles(adminInfo.ID)
	if err != nil {
		c.Error(conf.SERVER_ERROR, "查询失败: "+err.Error())
		return
	}

	// 返回登录信息
	c.Success(map[string]interface{}{
		"token": token,
		"user":  adminInfo.ToUserInfoRes(),
		"roles": roles,
	})
}

// GetUserInfo 获取当前登录管理员用户信息
// @Summary 获取当前登录管理员用户信息
// @Description 根据Token获取当前登录管理员用户的详细信息
// @Tags 后台-用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "{\"code\": 200, \"msg\": \"success\", \"data\": {...}}"
// @Failure 401 {object} map[string]interface{} "{\"code\": 401, \"msg\": \"未登录\"}"
// @router /api/admin/user/userinfo [get]
func (c *UserController) GetUserInfo() {
	// 返回用户信息
	userRoleModel := &adminModel.UserRole{}
	roles, err := userRoleModel.GetUserRoles(c.UserInfo.ID)
	if err != nil {
		c.Error(conf.SERVER_ERROR, "查询失败: "+err.Error())
		return
	}

	c.Success(map[string]interface{}{
		"user":  c.UserInfo.ToUserInfoRes(),
		"roles": roles,
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 已登录用户修改密码，需要验证旧密码，修改成功后first_login=0
// @Tags 后台-用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body ChangePasswordForm true "修改密码表单"
// @Success 200 {object} map[string]interface{} "{\"code\": 200, \"msg\": \"密码修改成功\"}"
// @Failure 400 {object} map[string]interface{} "{\"code\": 400, \"msg\": \"旧密码错误\"}"
// @router /api/admin/user/change-password [post]
func (c *UserController) ChangePassword() {
	var form admin.ChangePasswordForm
	err := c.ParseJson(&form)
	if err != nil {
		logs.Error("[UserController][ChangePassword] unmarshal error: %v", err)
		c.Error(conf.PARAMS_ERROR, "参数解析失败")
		return
	}

	// 参数校验
	if form.OldPassword == "" {
		c.Error(conf.PARAMS_ERROR, "旧密码不能为空")
		return
	}
	if form.NewPassword == "" {
		c.Error(conf.PARAMS_ERROR, "新密码不能为空")
		return
	}
	if len(form.NewPassword) < 6 {
		c.Error(conf.PARAMS_ERROR, "新密码长度不能少于6位")
		return
	}

	// 获取当前用户
	user := c.UserInfo
	if user == nil {
		c.Error(conf.UNAUTHORIZED, "用户未登录")
		return
	}

	// 修改密码（会验证旧密码并设置first_login=0）
	err = user.ChangePassword(form.OldPassword, form.NewPassword)
	if err != nil {
		logs.Error("[UserController][ChangePassword] change password error: %v", err)
		c.Error(conf.PARAMS_ERROR, "旧密码错误")
		return
	}

	c.Success(map[string]interface{}{
		"message": "密码修改成功",
	})
}

// Logout 管理员登出
// @Summary 管理员登出
// @Description 管理员登出，清除服务端会话（如果需要的话）
// @Tags 后台-用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "{\"code\": 200, \"msg\": \"登出成功\"}"
// @router /api/admin/user/logout [post]
func (c *UserController) Logout() {
	// JWT是无状态的，登出主要由前端处理（删除token）
	// 这里可以记录登出日志或做其他清理工作

	if c.UserInfo != nil {
		logs.Info("[UserController][Logout] User %s (ID: %d) logged out", c.UserInfo.Username, c.UserInfo.ID)
	}

	// token加入黑名单
	authHeader := c.Ctx.Input.Header("token")
	if authHeader == "" {
		// 尝试从 Authorization header 获取
		authHeader = c.Ctx.Input.Header("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}
	models.AddTokenToBlacklist(authHeader)

	c.Success(map[string]interface{}{
		"message": "登出成功",
	})
}
