package conf

// 标准 HTTP 错误码
const (
	SUCCESS      int64 = 200 // 成功
	PARAMS_ERROR int64 = 400 // 参数错误
	UNAUTHORIZED int64 = 401 // 未授权
	FORBIDDEN    int64 = 403 // 禁止访问
	NOT_FOUND    int64 = 404 // 未找到
	SERVER_ERROR int64 = 500 // 服务器错误
)

// 业务级错误码定义 (2000-2999)
const (
	USER_NOT_EXIST         = 2000 // 用户不存在
	USER_ALREADY_EXIST     = 2001 // 用户已存在
	PASSWORD_ERROR         = 2002 // 密码错误
	TOKEN_EXPIRED          = 2003 // Token已过期
	TOKEN_INVALID          = 2004 // Token无效
	LOGIN_REQUIRED         = 2005 // 需要登录
	VALID_CODE_EXPIRED     = 2006 // 验证码已过期
	VERIFY_CODE_ERROR      = 2007 // 验证码错误
	ACCOUNT_PASSWORD_ERROR = 2008 // 账号或密码错误
	DATA_NOT_FOUND         = 2009 // 数据不存在
)

// 用户注册/登录相关错误 (2100-2149)
const (
	ERROR_EMAIL_EMPTY              = 2100 // 邮箱不能为空
	ERROR_EMAIL_ALREADY_REGISTERED = 2101 // 邮箱已被注册
	ERROR_EMAIL_NOT_REGISTERED     = 2102 // 邮箱未注册
	ERROR_USERNAME_ALREADY_USED    = 2103 // 用户名已被使用
	ERROR_PASSWORD_STRENGTH        = 2104 // 密码强度不足
	ERROR_PAY_PASSWORD_FORMAT      = 2105 // 支付密码格式错误
	ERROR_INVITE_CODE_EMPTY        = 2106 // 邀请码不能为空
	ERROR_INVITE_CODE_NOT_EXIST    = 2107 // 邀请码不存在
	ERROR_VERIFY_CODE_INVALID      = 2108 // 验证码错误或已过期
	ERROR_SEND_CODE_FAILED         = 2109 // 发送验证码失败
	ERROR_REGISTER_FAILED          = 2110 // 注册失败
	ERROR_LOGIN_FAILED             = 2111 // 登录失败
	ERROR_USERNAME_PASSWORD_EMPTY  = 2112 // 用户名和密码不能为空
	ERROR_ACCOUNT_DISABLED         = 2113 // 账号已被禁用
	ERROR_USERNAME_PASSWORD_WRONG  = 2114 // 用户名或密码错误
	ERROR_REQUIRED_FIELDS_EMPTY    = 2115 // 必填字段不能为空
	ERROR_TYPE_INVALID             = 2116 // 类型参数错误
	ERROR_EMAIL_CODE_EMPTY         = 2117 // 邮箱和验证码不能为空
	ERROR_PASSWORD_TYPE_INVALID    = 2118 // 密码类型错误
	ERROR_NEW_PASSWORD_EMPTY       = 2119 // 新密码不能为空
	ERROR_RESET_PASSWORD_FAILED    = 2120 // 重置密码失败
	ERROR_SUBMIT_FAILED            = 2121 // 提交失败
	ERROR_GET_USER_INFO_FAILED     = 2122 // 获取用户信息失败
)

// 通用业务错误 2009-2099
const (
	ERROR_PARSE_FAILED     = 2009 // 参数解析失败
	ERROR_QUERY_FAILED     = 2010 // 查询失败
	ERROR_CREATE_FAILED    = 2011 // 创建失败
	ERROR_UPDATE_FAILED    = 2012 // 更新失败
	ERROR_DELETE_FAILED    = 2013 // 删除失败
	ERROR_RECORD_NOT_FOUND = 2014 // 记录不存在
	ERROR_RECORD_EXISTS    = 2015 // 记录已存在
	ERROR_NO_PERMISSION    = 2016 // 无权限操作
	ERROR_INVALID_ID       = 2017 // 无效的ID参数
	ERROR_MISSING_FIELDS   = 2018 // 缺少必填字段
	ERROR_INVALID_FORMAT   = 2019 // 格式错误
	ERROR_NO_UPDATE_FIELDS = 2020 // 没有可更新的字段
	ERROR_CHECK_FAILED     = 2021 // 检查失败
	ERROR_SAVE_FAILED      = 2022 // 保存失败
)
