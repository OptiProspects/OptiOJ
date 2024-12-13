package services

import (
	"OptiOJ/src/config"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, duration time.Duration) (string, error) {
	// 创建声明
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "OptiOJ",
		},
	}

	// 使用密钥创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获得完整的编码后的字符串令牌
	return token.SignedString(config.JWTSecret)
}

// ValidateRefreshToken 验证刷新令牌
func ValidateRefreshToken(refreshToken string) (uint, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWTSecret, nil
	})

	if err != nil {
		return 0, err
	}

	// 验证令牌是否有效
	if !token.Valid {
		return 0, jwt.ErrSignatureInvalid
	}

	// 获取声明
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, jwt.ErrInvalidType
	}

	// 验证令牌是否过期
	if claims.ExpiresAt.Before(time.Now()) {
		return 0, jwt.ErrTokenExpired
	}

	// 从Redis验证令牌是否在黑名单中
	blacklistKey := "token:blacklist:" + refreshToken
	exists, err := config.RedisClient.Exists(context.Background(), blacklistKey).Result()
	if err != nil {
		return 0, err
	}
	if exists == 1 {
		return 0, errors.New("token已被撤销")
	}

	return claims.UserID, nil
}

// 添加一个用于注销的函数
func InvalidateToken(token string, duration time.Duration) error {
	// 将令牌加入黑名单
	blacklistKey := "token:blacklist:" + token
	return config.RedisClient.Set(context.Background(), blacklistKey, "revoked", duration).Err()
}

// GenerateTokenPair 生成 Access Token 和 Refresh Token
func GenerateTokenPair(userID uint) (accessToken, refreshToken string, err error) {
	// 生成 Access Token
	accessToken, err = GenerateToken(userID, 2*time.Hour)
	if err != nil {
		return "", "", err
	}

	// 生成 Refresh Token
	refreshToken, err = GenerateToken(userID, 30*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateAccessToken 验证访问令牌
func ValidateAccessToken(accessToken string) (uint, error) {
	if accessToken == "" {
		return 0, errors.New("令牌为空")
	}

	// 解析令牌
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return config.JWTSecret, nil
	})

	if err != nil {
		return 0, fmt.Errorf("解析令牌失败: %v", err)
	}

	// 验证令牌是否有效
	if !token.Valid {
		return 0, jwt.ErrSignatureInvalid
	}

	// 获取声明
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, jwt.ErrInvalidType
	}

	// 验证令牌是否过期
	if claims.ExpiresAt.Before(time.Now()) {
		return 0, jwt.ErrTokenExpired
	}

	// 从Redis验证令牌是否在黑名单中
	blacklistKey := "token:blacklist:" + accessToken
	exists, err := config.RedisClient.Exists(context.Background(), blacklistKey).Result()
	if err != nil {
		return 0, fmt.Errorf("验证令牌黑名单失败: %v", err)
	}
	if exists == 1 {
		return 0, errors.New("令牌已被撤销")
	}

	// 验证令牌是否在Redis中存在（可选，用于实现服务器端控制会话）
	accessSessionKey := "access_token:" + accessToken
	userIDStr, err := config.RedisClient.Get(context.Background(), accessSessionKey).Result()
	if err != nil {
		return 0, fmt.Errorf("令牌已失效: %v", err)
	}

	// 验证Redis中存储的用户ID与令牌中的用户ID是否匹配
	redisUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("无效的用户ID: %v", err)
	}
	if uint(redisUserID) != claims.UserID {
		return 0, errors.New("令牌与会话不匹配")
	}

	return claims.UserID, nil
}
