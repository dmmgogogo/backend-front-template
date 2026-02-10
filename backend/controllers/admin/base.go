package admin

import (
	"e-woms/conf"
	"e-woms/models/admin"
	"e-woms/utils"
	"fmt"
	"slices"
	"std-library-slim/json"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

type BaseController struct {
	web.Controller
	i18n.Locale
	Code     int64
	Msg      string
	Result   interface{}
	UserId   int64
	UserInfo *admin.User
}

func (c *BaseController) Prepare() {
	c.Lang = c.Ctx.Input.Header("Language")
	if c.Lang == "" {
		c.Lang = "zh" //默认中文
	}

	// 排除登录注册接口
	path := c.Ctx.Request.URL.Path
	if slices.Contains(conf.NonLoginPathsAdmin, path) {
		return
	}

	// 开启IP白名单校验功能
	ip := c.Ctx.Input.IP()
	if !utils.IsIPInWhiteList(ip) {
		c.Error(conf.UNAUTHORIZED, "IP地址不在白名单内")
		return
	}

	// 从上下文获取user_id (由JWT中间件设置)
	userID, ok := c.Ctx.Input.GetData("user_id").(int64)
	if !ok {
		// 如果没有user_id,说明是未登录的公开接口,不需要加载UserInfo
		return
	}

	// 查询用户信息
	if userID > 0 {
		// 除非登录接口, 其他接口,都应该有值的, 不然权限接口肯定是过不去的
		c.UserId = userID

		// 查询用户信息
		var user admin.User
		err := user.GetByID(userID)
		if err != nil {
			logs.Error("[BaseController][Prepare] 查询用户失败: %v", err)
			c.Error(conf.UNAUTHORIZED)
			return
		}
		c.UserInfo = &user

		logs.Debug("[BaseController][Prepare] 查询用户成功: %v", json.String(c.UserInfo))
	}

	// 必须是有效玩家, 才能访问其他接口
	if c.UserInfo == nil {
		c.Error(conf.UNAUTHORIZED, "未登录")
		return
	}
}

// GetCurrentUserID 获取当前用户ID
func (c *BaseController) GetCurrentUserID() int64 {
	if c.UserInfo != nil {
		return c.UserInfo.ID
	}
	return 0
}

// TraceJson
func (c *BaseController) TraceJson() {
	res := map[string]interface{}{"code": c.Code, "msg": c.Msg, "data": c.Result}
	_ = c.JSONResp(res)
	c.StopRun()
}

func (c *BaseController) Success(obj interface{}) {
	c.Code = 200
	c.Result = obj
	c.TraceJson()
}

// 调用此函数，会终止当前函数内其他的逻辑处理（跳出）
func (c *BaseController) Error(code int64, msg ...string) {
	if code < 1 {
		code = 500
	}
	c.Code = code
	if c.Code > 0 {
		c.Msg = c.Tr("error." + fmt.Sprintf("%d", c.Code))
		// logs.Error(c.Code, c.Msg)
	}

	if len(msg) > 0 {
		c.Msg += "(" + strings.Join(msg, ",") + ")"
	}

	logs.Error("[Error] Code: %d, Message: %s", code, c.Msg)
	c.TraceJson()
}

func (c *BaseController) ParseJson(r interface{}) (err error) {
	err = json.ParseE(c.Ctx.Input.RequestBody, r)
	return
}

// GetAdminUserID 获取当前管理员用户ID
func (c *BaseController) GetAdminUserID() int64 {
	if c.UserInfo != nil {
		return c.UserInfo.ID
	}
	return 0
}

// GetAdminUsername 获取当前管理员用户名
func (c *BaseController) GetAdminUsername() string {
	if c.UserInfo != nil {
		return c.UserInfo.Username
	}
	return ""
}

// LogOperation 记录成功的操作日志
// operationType: create/update/delete/export
// module: 操作模块（如：客户管理、钱包管理）
// action: 具体操作描述（如：创建用户、导出钱包收益）
// targetType: 目标类型（如：user、wallet、order）
// targetID: 目标ID
// requestParams: 请求参数（业务相关）
func (c *BaseController) LogOperation(operationType, module, action, targetType string, targetID int64, requestParams map[string]interface{}) {
	// 异步记录日志（不阻塞主流程）
	go func() {
		err := admin.LogOperation(admin.LogOperationParams{
			AdminUserID:   c.GetAdminUserID(),
			AdminUsername: c.GetAdminUsername(),
			OperationType: operationType,
			Module:        module,
			Action:        action,
			TargetType:    targetType,
			TargetID:      targetID,
			RequestPath:   c.Ctx.Request.URL.Path,
			RequestMethod: c.Ctx.Request.Method,
			RequestParams: requestParams,
			IPAddress:     c.Ctx.Input.IP(),
			UserAgent:     c.Ctx.Request.UserAgent(),
			Status:        1, // 成功
			ErrorMsg:      "",
		})
		if err != nil {
			logs.Error("[LogOperation] Failed to log operation: %v", err)
		}
	}()
}

// LogOperationError 记录失败的操作日志
func (c *BaseController) LogOperationError(operationType, module, action string, errorMsg string) {
	go func() {
		_ = admin.LogOperation(admin.LogOperationParams{
			AdminUserID:   c.GetAdminUserID(),
			AdminUsername: c.GetAdminUsername(),
			OperationType: operationType,
			Module:        module,
			Action:        action,
			TargetType:    "",
			TargetID:      0,
			RequestPath:   c.Ctx.Request.URL.Path,
			RequestMethod: c.Ctx.Request.Method,
			RequestParams: nil,
			IPAddress:     c.Ctx.Input.IP(),
			UserAgent:     c.Ctx.Request.UserAgent(),
			Status:        0, // 失败
			ErrorMsg:      errorMsg,
		})
	}()
}
