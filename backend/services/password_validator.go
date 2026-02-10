package services

import (
	"regexp"
	"unicode"
)

// CheckPasswordStrength 检查密码强度
// 要求：至少8个字符，包含大写字母、小写字母、数字和特殊字符
func CheckPasswordStrength(password string) bool {
	// 至少8个字符
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	// 遍历每个字符检查类型
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 必须同时满足所有条件
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// IsDigit 检查字符串是否全部为数字
func IsDigit(s string) bool {
	if s == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^\d+$`, s)
	return matched
}
