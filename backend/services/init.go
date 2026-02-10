package services

import (
	"std-library-slim/dbase"
	"std-library-slim/json"
	"std-library-slim/redis"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

func InitMysql() {
	opt := dbase.Opt{}
	err := json.ParseE(web.AppConfig.DefaultString("MYSQL_CONFIG", ""), &opt)
	if err != nil {
		logs.Error("Failed to init MySQL: %v", err)
		panic(err)
	}
	dbase.Init(&opt)
}

// InitRedis 初始化Redis
func InitRedis() {
	opt := redis.Opt{}
	err := json.ParseE(web.AppConfig.DefaultString("REDIS_CONFIG", ""), &opt)
	if err != nil {
		logs.Error("Failed to init Redis: %v", err)
		panic(err)
	}
	redis.Init(&opt)
}
