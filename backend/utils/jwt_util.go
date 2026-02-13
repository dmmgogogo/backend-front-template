package utils

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v4"
)

// getJWTSecret 从 conf 读取 JWT_SECRET，未配置时返回错误
func getJWTSecret() ([]byte, error) {
	s, err := web.AppConfig.String("JWT_SECRET")
	if err != nil || s == "" {
		return nil, errors.New("JWT_SECRET not configured")
	}
	return []byte(s), nil
}

// Claims JWT 声明
type Claims struct {
	UserID   int64  `json:"user_id"`
	DeviceID string `json:"device_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(userID int64, deviceID string) (string, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// GenerateRefreshToken 生成刷新 token
func GenerateRefreshToken(userID int64, deviceID string) (string, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 天
	claims := &Claims{
		UserID:   userID,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证 token
func ValidateToken(tokenString string) (bool, *Claims) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return false, nil
	}
	return true, claims
}
