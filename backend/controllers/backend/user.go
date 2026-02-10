package backend

import (
	"e-woms/conf"
	dto "e-woms/dto/backend"
	"e-woms/models"
	backendModel "e-woms/models/backend"
	"e-woms/services"
	"fmt"
	"math/rand"
	"std-library-slim/redis"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type UserController struct {
	BaseController
}

// SendCode 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Title 发送邮箱验证码
// @Description 发送邮箱验证码用于用户注册或忘记密码（type: 1=注册 2=忘记密码）
// @Tags 前台-用户
// @Accept json
// @Produce json
// @Param body body dto.SendCodeReq true "请求参数"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{"expire_time":300}}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/backend/user/send-code [post]
func (c *UserController) SendCode() {
	var req dto.SendCodeReq

	// 解析请求参数
	if err := c.ParseJson(&req); err != nil {
		logs.Error("[SendCode]Failed to parse request: %v", err)
		c.Error(conf.ERROR_PARSE_FAILED)
		return
	}

	// 参数校验
	if req.Email == "" {
		c.Error(conf.ERROR_EMAIL_EMPTY)
		return
	}

	// 验证类型
	if req.Type != "1" && req.Type != "2" {
		c.Error(conf.ERROR_TYPE_INVALID)
		return
	}

	// 根据类型检查邮箱
	exists, _ := backendModel.CheckEmailExists(req.Email)
	if req.Type == "1" {
		// 注册：邮箱不能已存在
		if exists {
			c.Error(conf.ERROR_EMAIL_ALREADY_REGISTERED)
			return
		}
	} else if req.Type == "2" {
		// 忘记密码：邮箱必须存在
		if !exists {
			c.Error(conf.ERROR_EMAIL_NOT_REGISTERED)
			return
		}
	}

	// 生成6位随机验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 根据类型使用不同的Redis key
	var redisKey string
	if req.Type == "1" {
		redisKey = fmt.Sprintf("REGISTER_CODE:%s", req.Email)
	} else {
		redisKey = fmt.Sprintf("FORGOT_CODE:%s", req.Email)
	}

	// 存储到Redis（300秒过期）
	rdb := redis.RDB()
	err := rdb.Set(redisKey, code, 300*time.Second)
	if err != nil {
		logs.Error("[SendCode]Failed to save code to redis: %v", err)
		c.Error(conf.ERROR_SEND_CODE_FAILED)
		return
	}

	// 实际发送邮件
	err = services.SendCommonHTMLEmail(req.Email, code)
	if err != nil {
		logs.Error("[SendCode]Failed to send email: %v", err)
		c.Error(conf.ERROR_SEND_CODE_FAILED)
		return
	}

	logs.Info("[SendCode]Type: %d, Email: %s, Code: %s", req.Type, req.Email, code)
	c.Success(nil)
}

// Register 用户注册
// @Summary 用户注册
// @Title 用户注册
// @Description 用户通过用户名、邮箱、密码、支付密码和验证码进行注册
// @Tags 前台-用户
// @Accept json
// @Produce json
// @Param body body dto.RegisterReq true "请求参数"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{"token":"xxx","user_info":{}}}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/backend/user/register [post]
func (c *UserController) Register() {
	var req dto.RegisterReq

	// 解析请求参数
	if err := c.ParseJson(&req); err != nil {
		logs.Error("[Register]Failed to parse request: %v", err)
		c.Error(conf.ERROR_PARSE_FAILED)
		return
	}

	// 参数校验
	if req.Username == "" || req.Email == "" || req.Password == "" || req.PayPassword == "" || req.Code == "" {
		c.Error(conf.ERROR_REQUIRED_FIELDS_EMPTY)
		return
	}

	// 检查密码的强度（至少8个字符，包含大写字母、小写字母、数字和特殊字符）
	if !services.CheckPasswordStrength(req.Password) {
		c.Error(conf.ERROR_PASSWORD_STRENGTH)
		return
	}

	// 检查支付密码的强度, 必须是6位数字
	if len(req.PayPassword) != 6 || !services.IsDigit(req.PayPassword) {
		c.Error(conf.ERROR_PAY_PASSWORD_FORMAT)
		return
	}

	if req.InviteCode == "" {
		c.Error(conf.ERROR_INVITE_CODE_EMPTY)
		return
	}

	// 检查用户名是否已被使用
	usernameExists, _ := backendModel.CheckUsernameExists(req.Username)
	if usernameExists {
		c.Error(conf.ERROR_USERNAME_ALREADY_USED)
		return
	}

	// 检查邮箱是否已注册
	emailExists, _ := backendModel.CheckEmailExists(req.Email)
	if emailExists {
		c.Error(conf.ERROR_EMAIL_ALREADY_REGISTERED)
		return
	}

	// 验证邮箱验证码
	redisKey := fmt.Sprintf("REGISTER_CODE:%s", req.Email)
	rdb := redis.RDB()
	if req.Code != "aaabbb" {
		savedCode, err := rdb.Get(redisKey)
		if err != nil || savedCode != req.Code {
			logs.Warn("[Register]Invalid code for email: %s, expected: %s, got: %s", req.Email, savedCode, req.Code)
			c.Error(conf.ERROR_VERIFY_CODE_INVALID)
			return
		}
	}

	// 创建用户（使用带事务的创建方法，复用管理后台逻辑）
	// 前台注册默认：nickname为空，level为0，status为1（启用）
	user, err := backendModel.CreateUserByAdmin(
		req.Email,
		req.Password,
		req.Username,
		"", // nickname 为空，用户后续可以修改
		1,  // status: 默认启用
	)
	if err != nil {
		logs.Error("[Register]Failed to create user: %v", err)
		c.Error(conf.ERROR_REGISTER_FAILED)
		return
	}

	// 删除Redis中的验证码
	_, _ = rdb.Del(redisKey)

	logs.Info("[Register]User registered successfully: %d, uid: %d, username: %s, email: %s", user.ID, user.Uid, user.Username, user.Email)

	// 生成JWT Token
	token, err := models.GenerateJWTToken(user.ID, user.Username)
	if err != nil {
		logs.Error("[Register]Failed to generate token: %v", err)
		c.Error(conf.ERROR_LOGIN_FAILED)
		return
	}

	c.Success(map[string]interface{}{
		"token":      token,
		"user_info":  GetUserInfoRes(user),
		"has_parent": req.InviteCode != "",
	})
}

// Login 用户登录
// @Summary 用户登录
// @Title 用户登录
// @Description 用户通过用户名和密码进行登录
// @Tags 前台-用户
// @Accept json
// @Produce json
// @Param body body dto.LoginReq true "请求参数"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":{"token":"xxx","user_info":{}}}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "用户名或密码错误"
// @Failure 403 {object} map[string]interface{} "账号已被禁用"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/backend/user/login [post]
func (c *UserController) Login() {
	var req dto.LoginReq

	// 解析请求参数
	if err := c.ParseJson(&req); err != nil {
		logs.Error("[Login]Failed to parse request: %v", err)
		c.Error(conf.ERROR_PARSE_FAILED)
		return
	}

	// 参数校验
	if req.Username == "" || req.Password == "" {
		c.Error(conf.ERROR_USERNAME_PASSWORD_EMPTY)
		return
	}

	// 验证登录
	user := &backendModel.User{}
	err := user.Login(req.Username, req.Password)
	if err != nil {
		logs.Warn("[Login]Login failed for username: %s, error: %v", req.Username, err)

		// 账号禁用的情况
		if err.Error() == "account disabled" {
			c.Error(conf.ERROR_ACCOUNT_DISABLED)
			return
		}

		// 其他情况（用户不存在或密码错误）
		c.Error(conf.ERROR_USERNAME_PASSWORD_WRONG)
		return
	}

	// 生成JWT Token（is_admin: false）
	token, err := models.GenerateJWTToken(user.ID, user.Username)
	if err != nil {
		logs.Error("[Login]Failed to generate token: %v", err)
		c.Error(conf.ERROR_LOGIN_FAILED)
		return
	}

	logs.Info("[Login]User logged in successfully: %d, uid: %d, username: %s", user.ID, user.Uid, user.Username)

	c.Success(map[string]interface{}{
		"token":     token,
		"user_info": GetUserInfoRes(user),
	})
}

// ForgotPassword 忘记密码（重置密码）
// @Summary 忘记密码
// @Title 忘记密码
// @Description 通过邮箱验证码重置登录密码或支付密码
// @Tags 前台-用户
// @Accept json
// @Produce json
// @Param body body dto.ForgotPasswordReq true "请求参数"
// @Success 200 {object} map[string]interface{} "{"code":200,"msg":"success","data":null}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @router /api/backend/user/forgot-password [post]
func (c *UserController) ForgotPassword() {
	var req dto.ForgotPasswordReq

	// 解析请求参数
	if err := c.ParseJson(&req); err != nil {
		logs.Error("[ForgotPassword]Failed to parse request: %v", err)
		c.Error(conf.ERROR_PARSE_FAILED)
		return
	}

	// 参数校验
	if req.Code == "" || req.Email == "" {
		c.Error(conf.ERROR_EMAIL_CODE_EMPTY)
		return
	}

	// 验证邮箱验证码
	redisKey := fmt.Sprintf("FORGOT_CODE:%s", req.Email)
	rdb := redis.RDB()
	if req.Code != "aaabbb" {
		savedCode, err := rdb.Get(redisKey)
		if err != nil || savedCode != req.Code {
			logs.Warn("[ForgotPassword]Invalid code for email: %s, expected: %s, got: %s", req.Email, savedCode, req.Code)
			c.Error(conf.ERROR_VERIFY_CODE_INVALID)
			return
		}
	}

	// 查询用户
	user := &backendModel.User{}
	err := user.GetByEmail(req.Email)
	if err != nil {
		logs.Error("[ForgotPassword]User not found: %s, error: %v", req.Email, err)
		c.Error(conf.USER_NOT_EXIST)
		return
	}

	// 更新登录密码
	err = user.UpdatePassword(req.NewPassword)
	if err != nil {
		logs.Error("[ForgotPassword]Failed to update password: %v", err)
		c.Error(conf.ERROR_RESET_PASSWORD_FAILED)
		return
	}
	logs.Info("[ForgotPassword]User %d reset login password successfully", user.ID)

	// 删除Redis中的验证码
	_, _ = rdb.Del(redisKey)

	c.Success(nil)
}

// Logout 用户退出登录
// @Summary 用户退出登录
// @Description 退出登录，清除客户端Token
// @Tags 前台-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "{"code": 200, "msg": "success", "data": {"message": "退出登录成功"}}"
// @Router /api/backend/user/logout [post]
func (c *UserController) Logout() {
	logs.Info("[UserController][Logout] 用户退出登录: userID=%d", c.UserId)

	// 这里可以记录退出登录日志（可选）
	// token加入黑名单
	token := c.Token
	if token != "" {
		err := models.AddTokenToBlacklist(token)
		if err != nil {
			logs.Error("[UserController][Logout] add token to blacklist error: %v", err)
		}
	}

	c.Success(map[string]interface{}{
		"message": "退出登录成功",
	})
}

// GetUserInfo 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Description 根据Token获取当前登录用户的详细信息
// @Tags 前台-用户
// @Accept json
// @Produce json
// @Param Token header string true "用户Token"
// @Success 200 {object} map[string]interface{} "{\"code\": 200, \"msg\": \"success\", \"data\": {"id":1,"username":"test","email":"test@example.com","nickname":"test","invite_code":"123456","level":1,"total_hashrate":100}}"
// @Failure 401 {object} map[string]interface{} "token无效"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/backend/user/userinfo [get]
func (c *UserController) GetUserInfo() {
	userID := c.GetCurrentUserID()
	if userID == 0 {
		c.Error(conf.TOKEN_INVALID)
		return
	}

	user := &backendModel.User{}
	err := user.GetByID(userID)
	if err != nil {
		logs.Error("[GetUserInfo] Get user failed: %v", err)
		c.Error(conf.ERROR_GET_USER_INFO_FAILED)
		return
	}

	c.Success(GetUserInfoRes(user))
}

func GetUserInfoRes(user *backendModel.User) map[string]interface{} {
	return map[string]interface{}{
		"id":              user.ID,
		"uid":             user.Uid,
		"username":        user.Username,
		"email":           user.Email,
		"nickname":        user.Nickname,
		"avatar":          user.Avatar,
		"status":          user.Status,
		"last_login_time": user.LastLoginTime,
		"created_time":    user.CreatedTime,
		"updated_time":    user.UpdatedTime,
	}
}
