package conf

const (
	// 万能token
	TOKEN_KING                      = "fNY0TlRVYNuf9emiTdD5fXralZmzF6"
	KeyPhoneValidCode               = "KeyPhoneValidCode:%v%v"     //手机验证码
	KeyPhoneValidCodeLock           = "KeyPhoneValidCodeLock:%v%v" //手机验证码时间锁
	KeyPhoneValidCodeExpireTime     = 60 * 5
	KeyPhoneValidCodeLockExpireTime = 10

	// 邮箱验证码
	KeyEmailValidCode               = "KeyEmailValidCode:%v"     //邮箱验证码
	KeyEmailValidCodeLock           = "KeyEmailValidCodeLock:%v" //邮箱验证码时间锁
	KeyEmailValidCodeExpireTime     = 60 * 5
	KeyEmailValidCodeLockExpireTime = 10
)
