package backend

import (
	"e-woms/conf"
	"e-woms/utils"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

type BaseController struct {
	web.Controller
	i18n.Locale
	Code   int64
	Msg    string
	Result interface{}
	Token  string
	UserId int64
	// UserInfo *apiModel.User
}

func (c *BaseController) Prepare() {
	// accept-language = zh-CN, zh-SG, zh-HK, en-US, en-GB 等，默认zh
	// ⚠️ 必须先设置语言，否则 Error() 方法中的 c.Tr() 无法正确翻译
	acceptLanguage := c.Ctx.Input.Header("Accept-Language")
	c.Lang = "zh" // 默认中文
	if acceptLanguage != "" {
		// 取第一个语言（优先级最高）
		firstLang := strings.Split(acceptLanguage, ",")[0]
		// 去除权重标识（如 ;q=0.9）
		firstLang = strings.Split(firstLang, ";")[0]
		firstLang = strings.TrimSpace(firstLang)

		// 判断语言类型：所有 zh-* 归为 zh，所有 en-* 归为 en
		if strings.HasPrefix(firstLang, "zh") {
			c.Lang = "zh"
		} else if strings.HasPrefix(firstLang, "en") {
			c.Lang = "en"
		}
	}
	logs.Debug("[BaseController][Prepare] acceptLanguage: %s, Lang: %s", acceptLanguage, c.Lang)

	// 排除登录注册接口
	path := c.Ctx.Request.URL.Path
	if slices.Contains(conf.NonLoginPathsBackend, path) {
		return
	}

	// 解析 JWT Token
	authHeader := c.Ctx.Input.Header("Authorization")
	if authHeader == "" {
		logs.Error("[BaseController] 缺少 Authorization 请求头")
		c.Error(conf.UNAUTHORIZED, c.Tr("api.unauthorized"))
		return
	}

	// 验证 Bearer 格式
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		logs.Error("[BaseController] Authorization 格式错误")
		c.Error(conf.UNAUTHORIZED, c.Tr("api.invalid_token"))
		return
	}

	// 解析 Token
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		logs.Error("[BaseController] Token 解析失败: %v", err)
		c.Error(conf.TOKEN_INVALID, c.Tr("api.token_invalid"))
		return
	}

	// 设置用户 ID
	c.UserId = claims.UserID
	logs.Debug("[BaseController] 解析 Token 成功, UserID: %d", c.UserId)
}

// GetCurrentUserID 获取当前用户ID
func (c *BaseController) GetCurrentUserID() int64 {
	// if c.UserInfo != nil {
	// 	return c.UserInfo.ID
	// }
	return c.UserId
}

// TraceJson
func (c *BaseController) TraceJson() {
	res := map[string]interface{}{"code": c.Code, "msg": c.Msg, "data": c.Result}
	c.Data["json"] = res
	c.ServeJSON()
	c.StopRun()
}

func (c *BaseController) Success(obj interface{}) {
	c.Code = 200
	c.Msg = "success"
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

// ParseJson 解析 JSON 请求体
func (c *BaseController) ParseJson(r interface{}) error {
	err := json.Unmarshal(c.Ctx.Input.RequestBody, r)
	if err != nil {
		logs.Error("[ParseJson]Failed to parse: %v", err)
	}
	return err
}
