package utils

import (
	"fmt"
	"time"

	redisLib "std-library-slim/redis"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

const (
	// Redis key 前缀
	IPWhitelistRedisKey = "ip_whitelist"
)

// IsIPWhitelistEnabled 检查 IP 白名单功能是否开启
func IsIPWhitelistEnabled() bool {
	enabled, err := web.AppConfig.Bool("IP_WHITELIST_ENABLED")
	if err != nil {
		logs.Debug("[IP Whitelist] IP_WHITELIST_ENABLED not configured, default to false")
		return false
	}
	return enabled
}

// IsIPInWhiteList 检查 IP 是否在白名单中
func IsIPInWhiteList(ip string) bool {
	// 如果白名单功能未开启，直接返回 true
	if !IsIPWhitelistEnabled() {
		return true
	}

	logs.Debug("[IP Whitelist] Checking IP %s in whitelist", ip)

	// 从 Redis 检查 IP 是否在白名单中
	rdb := redisLib.RDB()

	isMember, err := rdb.SIsMember(IPWhitelistRedisKey, ip)
	if err != nil {
		logs.Error("[IP Whitelist] Failed to check IP in whitelist: %v", err)
		// 如果 Redis 出错，为了安全起见，拒绝访问
		return false
	}

	if !isMember {
		logs.Warn("[IP Whitelist] IP %s not in whitelist", ip)
	}

	return isMember
}

// AddIPToWhitelist 添加 IP 到白名单
func AddIPToWhitelist(ip string) error {
	rdb := redisLib.RDB()

	// 使用 Redis Set 存储白名单
	_, err := rdb.SAdd(IPWhitelistRedisKey, ip)
	if err != nil {
		logs.Error("[IP Whitelist] Failed to add IP %s to whitelist: %v", ip, err)
		return fmt.Errorf("添加 IP 到白名单失败: %v", err)
	}

	logs.Info("[IP Whitelist] Successfully added IP %s to whitelist", ip)
	return nil
}

// RemoveIPFromWhitelist 从白名单移除 IP
func RemoveIPFromWhitelist(ip string) error {
	rdb := redisLib.RDB()

	_, err := rdb.SRem(IPWhitelistRedisKey, ip)
	if err != nil {
		logs.Error("[IP Whitelist] Failed to remove IP %s from whitelist: %v", ip, err)
		return fmt.Errorf("从白名单移除 IP 失败: %v", err)
	}

	logs.Info("[IP Whitelist] Successfully removed IP %s from whitelist", ip)
	return nil
}

// GetAllWhitelistIPs 获取所有白名单 IP
func GetAllWhitelistIPs() ([]string, error) {
	rdb := redisLib.RDB()

	ips, err := rdb.SMembers(IPWhitelistRedisKey)
	if err != nil {
		logs.Error("[IP Whitelist] Failed to get all whitelist IPs: %v", err)
		return nil, fmt.Errorf("获取白名单失败: %v", err)
	}

	return ips, nil
}

// CountWhitelistIPs 获取白名单 IP 数量
func CountWhitelistIPs() (int64, error) {
	rdb := redisLib.RDB()

	count, err := rdb.SCard(IPWhitelistRedisKey)
	if err != nil {
		logs.Error("[IP Whitelist] Failed to count whitelist IPs: %v", err)
		return 0, fmt.Errorf("获取白名单数量失败: %v", err)
	}

	return count, nil
}

// GetIPWhitelistManageKey 获取 IP 白名单管理密钥
func GetIPWhitelistManageKey() string {
	key, err := web.AppConfig.String("IP_WHITELIST_MANAGE_KEY")
	if err != nil || key == "" {
		logs.Warn("[IP Whitelist] IP_WHITELIST_MANAGE_KEY not configured")
		return ""
	}
	return key
}

// LogIPWhitelistOperation 记录 IP 白名单操作日志
func LogIPWhitelistOperation(action, ip, clientIP string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	logs.Info("[IP Whitelist Operation] Action: %s, IP: %s, ClientIP: %s, Status: %s, Time: %s",
		action, ip, clientIP, status, time.Now().Format("2006-01-02 15:04:05"))
}
