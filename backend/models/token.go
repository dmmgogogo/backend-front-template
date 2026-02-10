package models

import (
	"fmt"
	"time"

	"std-library-slim/redis"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
)

const TOKEN_BLACKLIST_PREFIX = "token_blacklist:"

// AddTokenToBlacklist 将token加入黑名单
func AddTokenToBlacklist(tokenString string) error {
	// 解析token获取过期时间
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(web.AppConfig.DefaultString("JWT_SECRET", "")), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	// 获取token的过期时间
	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("invalid token expiration")
	}

	// 计算剩余有效期
	expTime := time.Unix(int64(exp), 0)
	ttl := time.Until(expTime)
	if ttl <= 0 {
		return nil // token已过期，无需加入黑名单
	}

	// 将token加入Redis黑名单，使用剩余有效期作为过期时间
	key := TOKEN_BLACKLIST_PREFIX + tokenString
	err = redis.RDB().Set(key, "1", ttl)
	if err != nil {
		logs.Error("[AddTokenToBlacklist]Failed to add token to blacklist: %v", err)
		return err
	}

	return nil
}

// IsTokenBlacklisted 检查token是否在黑名单中
func IsTokenBlacklisted(tokenString string) bool {
	key := TOKEN_BLACKLIST_PREFIX + tokenString
	exists, err := redis.RDB().Exists(key)
	if err != nil {
		logs.Error("[IsTokenBlacklisted]Failed to check token blacklist: %v", err)
		return false
	}
	return exists
}
