package admin

import "github.com/golang-jwt/jwt/v4"

// LoginForm 登录表单（支持邮箱登录）
type LoginForm struct {
	Username   string `json:"username"`    // 用户名（必填）
	Password   string `json:"password"`    // 密码（必填）
	VerifyCode string `json:"verify_code"` // Google验证码（必填）
}

// Claims JWT Claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  int    `json:"is_admin"`
	jwt.RegisteredClaims
}

// ChangePasswordForm 修改密码表单
type ChangePasswordForm struct {
	OldPassword string `json:"old_password"` // 旧密码
	NewPassword string `json:"new_password"` // 新密码
}

// ForgotPasswordForm 忘记密码表单
type ForgotPasswordForm struct {
	Email string `json:"email"` // 邮箱
}

// ================ 客户管理相关 DTO ================

// CustomerListReq 客户列表查询请求
type CustomerListReq struct {
	Page       int    `json:"page"`        // 页码
	PageSize   int    `json:"page_size"`   // 每页数量
	Email      string `json:"email"`       // 邮箱（模糊搜索）
	InviteCode string `json:"invite_code"` // 邀请码（模糊搜索）
	Status     int    `json:"status"`      // 状态：-1=全部 0=禁用 1=正常
}

// CustomerCreateReq 创建客户请求
type CustomerCreateReq struct {
	Email    string `json:"email"`    // 邮箱（必填，唯一）
	Password string `json:"password"` // 登录密码（必填）
	Username string `json:"username"` // 用户名（可选）
	Nickname string `json:"nickname"` // 昵称（可选）
	Status   int    `json:"status"`   // 状态（可选，默认1）
}

// CustomerResetPasswordReq 重置客户密码请求
type CustomerResetPasswordReq struct {
	NewPassword  string `json:"new_password"`  // 新密码
	PasswordType int    `json:"password_type"` // 密码类型：1=登录密码 2=支付密码
}

// CustomerRes 客户响应结构
type CustomerRes struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Nickname      string `json:"nickname"`
	Avatar        string `json:"avatar"`
	Status        int    `json:"status"`
	LastLoginTime int64  `json:"last_login_time"`
	CreatedTime   int64  `json:"created_time"`
	UpdatedTime   int64  `json:"updated_time"`
}
