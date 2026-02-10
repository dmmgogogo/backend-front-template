package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	mrand "math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
)

// 数组元素去重
func RemoveDuplicates(input []string) []string {
	// 用 map 记录已经遇到的字符串
	seen := make(map[string]bool)
	var result []string

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

// GenerateJWTToken 生成JWT token（企业用户）
func GenerateJWTToken(userId int64, username string) (string, error) {
	jwtSecret, err := web.AppConfig.String("JWT_SECRET")
	if jwtSecret == "" || err != nil {
		return "", errors.New("failed to get JWT secret")
	}

	// 默认365天
	jwtExpireTime := web.AppConfig.DefaultInt("JWT_SECRET_EXPIRE_TIME", 365*24)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userId
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(jwtExpireTime)).Unix()

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenerateAdminJWTToken 生成管理员JWT token（包含is_admin标识）
func GenerateAdminJWTToken(userId int64, username string) (string, error) {
	jwtSecret, err := web.AppConfig.String("JWT_SECRET")
	if jwtSecret == "" || err != nil {
		return "", errors.New("failed to get JWT secret")
	}

	// 默认365天
	jwtExpireTime := web.AppConfig.DefaultInt("JWT_SECRET_EXPIRE_TIME", 365*24)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userId
	claims["username"] = username
	claims["is_admin"] = true // 是管理后台用户
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(jwtExpireTime)).Unix()

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseJWTToken 解析JWT token
func ParseJWTToken(tokenString string) (map[string]interface{}, error) {
	jwtSecret, err := web.AppConfig.String("JWT_SECRET")
	if err != nil {
		return nil, errors.New("failed to get JWT secret")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 定义响应结构体
type IPInfo struct {
	Country string `json:"country"`
}

// 根据IP获取当前国家，目前不支持IPv6
func GetDomesticIP(ip string) (country string, err error) {
	// 使用ipinfo.io API 查询IP的地理位置
	url := fmt.Sprintf("http://ipinfo.io/%s/json", ip)
	resp, err := http.Get(url)
	if err != nil {
		logs.Error("Failed to get ipinfo.io info from ipinfo.io", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("Failed to read ipinfo.io info from ipinfo.io", err)
		return
	}

	// 解析JSON响应
	var ipInfo IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		logs.Error("Failed to unmarshal ipinfo.io info from ipinfo.io", err)
		return
	}

	logs.Debug("ipinfo.io info from ipinfo.io is", ipInfo)

	return ipInfo.Country, nil
}

// 直接发文本
func SendToTelegramOnlyText(content string) {
	go func() {
		// 从配置文件获取 token 和 group_id
		token, err := web.AppConfig.String("telegram::token")
		if err != nil {
			logs.Error("[Telegram] 获取token失败: %v", err)
			return
		}

		groupID, err := web.AppConfig.String("telegram::group_id")
		if err != nil {
			logs.Error("[Telegram] 获取group_id失败: %v", err)
			return
		}

		// 构建 Telegram API URL
		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

		// 构建请求参数
		params := url.Values{}
		params.Add("chat_id", groupID)
		params.Add("text", TruncateForTelegram(content))
		// params.Add("parse_mode", "HTML") // 支持HTML格式

		// 发送请求
		resp, err := http.PostForm(apiURL, params)
		if err != nil {
			logs.Error("[Telegram] 发送消息失败: %v", err)
			return
		}
		defer resp.Body.Close()

		// 检查响应状态
		if resp.StatusCode != http.StatusOK {
			logs.Error("[Telegram] 发送消息失败，状态码: %d", resp.StatusCode)
			return
		}

		logs.Info("[Telegram] 消息发送成功: %s", content)
	}()
}

// SendToTelegram 发送消息到Telegram群组
func SendToTelegram(content string) {
	if web.AppConfig.DefaultBool("LOCAL_CLOSE_TELEGRAM", false) {
		return
	}

	go func() {
		logs.Debug("[Telegram] 发送消息内容: content:[%s]", content)

		SendToTelegramOnlyText(content)
	}()
}

// Telegram消息长度限制
const TelegramMaxLength = 4096

// 截取单条消息
func TruncateForTelegram(message string) string {
	runes := []rune(message)
	if len(runes) <= TelegramMaxLength {
		return message
	}
	return string(runes[:TelegramMaxLength-3]) + "..."
}

// RandomNick 用户+随机字母3位+随机数字8位
func RandomNick() string {
	mrand.Seed(time.Now().UnixNano())
	nick := make([]byte, 3)
	for i := range nick {
		nick[i] = byte(mrand.Intn(26) + 97) // 生成小写字母
	}

	// 生成8位随机数字
	randNum := mrand.Intn(90000000) + 10000000 // 生成8位数字(10000000-99999999)

	return "用户" + string(nick) + strconv.Itoa(randNum)
}

// Bytes2String Bytes2String 将byte数组，作为普通数组转换为字符串，而非按 ASCII 码转换
// []byte{97,98} => "9798"
func Bytes2String(bs []byte) (s string) {
	for _, b := range bs {
		s += fmt.Sprint(b)
	}
	return
}
