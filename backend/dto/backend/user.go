package dto

// SendCodeReq 发送验证码请求
type SendCodeReq struct {
	Email string `json:"email"` // 邮箱（必填）
	Type  string `json:"type"`  // 1:注册 2:忘记密码
}

// RegisterReq 用户注册请求
type RegisterReq struct {
	Username    string `json:"username"`     // 用户名（必填）
	Email       string `json:"email"`        // 邮箱（必填）
	Password    string `json:"password"`     // 密码（必填）
	PayPassword string `json:"pay_password"` // 支付密码（必填）
	Code        string `json:"code"`         // 验证码（必填）
	InviteCode  string `json:"invite_code"`  // 邀请码（可选）
}

// ForgotPasswordReq 忘记密码请求
type ForgotPasswordReq struct {
	Email        string `json:"email"`         // 邮箱（必填）
	Code         string `json:"code"`          // 验证码（必填）
	NewPassword  string `json:"new_password"`  // 新密码（password_type=1时必填）
	PayPassword  string `json:"pay_password"`  // 新支付密码（password_type=2时必填）
	PasswordType int    `json:"password_type"` // 1:登录密码 2:支付密码
}

// LoginReq 用户登录请求
type LoginReq struct {
	Username string `json:"username"` // 用户名（必填）
	Password string `json:"password"` // 密码（必填）
}
