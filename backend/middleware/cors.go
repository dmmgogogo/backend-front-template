package middleware

import (
	"github.com/beego/beego/v2/server/web/context"
)

func Cors(ctx *context.Context) {
	// 设置基本的CORS响应头
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Headers", "Origin,Authorization,Content-Type,Accept,X-Requested-With,sec-ch-ua,sec-ch-ua-platform,sec-ch-ua-mobile,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Token,token")
	ctx.Output.Header("Access-Control-Expose-Headers", "Authorization,Content-Length,Access-Control-Allow-Origin")
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")

	// 只处理OPTIONS请求
	if ctx.Input.Method() == "OPTIONS" {
		ctx.Output.SetStatus(200)
		ctx.ResponseWriter.Write([]byte(""))
		return
	}
}
