package models

import (
	"time"

	"std-library-slim/redis"
	"yzyw/utils"

	"github.com/beego/beego/v2/core/logs"
)

const TOKEN_BLACKLIST_PREFIX = "token_blacklist:"

// AddTokenToBlacklist 将token加入黑名单
func AddTokenToBlacklist(tokenString string) error {
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		return err
	}

	if claims.ExpiresAt == nil {
		return nil // 无过期时间，不加入黑名单
	}
	ttl := time.Until(claims.ExpiresAt.Time)
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
