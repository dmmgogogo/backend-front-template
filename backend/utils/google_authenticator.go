package utils

import (
	"github.com/pquerna/otp/totp"
)

// VerifyGoogleAuthCode 验证 Google 验证码
// secret: 用户的 Google Authenticator Secret Key (存储在数据库中)
// code: 用户输入的 6 位验证码
// 返回: true=验证通过, false=验证失败
func VerifyGoogleAuthCode(secret, code string) bool {
	return totp.Validate(code, secret)
}
