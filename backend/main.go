package main

import (
	"e-woms/routers"
	"e-woms/services"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

func main() {
	// 初始化
	initLogger()
	initLocales()
	// 初始化路由
	routers.InitRouters()
	// 初始化mysql
	services.InitMysql()
	// 初始化redis
	services.InitRedis()

	logs.Debug("Starting server... | version: v1.0.11")
	web.Run()
}

// 初始化多语言
func initLocales() {
	langs := []string{"zh", "en"}
	for _, lang := range langs {
		logs.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			logs.Error("Fail to set message file: " + err.Error())
			return
		}
	}
}

// 初始化log
func initLogger() {
	logs.SetLogger(logs.AdapterFile, `{ "filename": "logs/app.log", "daily": true, "maxdays": 7}`)

	// logs.SetLogger(logs.AdapterConsole)
	logs.SetLevel(logs.LevelDebug) // 设置日志级别为 Debug
	logs.SetLogFuncCall(true)      // 显示文件名和行号
	logs.SetLogFuncCallDepth(3)
}
