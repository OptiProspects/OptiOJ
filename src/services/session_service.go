package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"OptiOJ/src/config"
)

type DeviceInfo struct {
	UserAgent   string    `json:"user_agent"`
	IP          string    `json:"ip"`
	LastActive  time.Time `json:"last_active"`
	LastRefresh time.Time `json:"last_refresh"`
}

type SessionInfo struct {
	SessionID  string     `json:"session_id"` // 用于标识会话的唯一ID
	DeviceInfo DeviceInfo `json:"device_info"`
	CreatedAt  time.Time  `json:"created_at"`
}

// 生成会话ID
func generateSessionID(refreshToken string) string {
	hash := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(hash[:])
}

// GetActiveSessions 获取用户的所有活跃会话
func GetActiveSessions(userID uint) ([]SessionInfo, error) {
	ctx := context.Background()
	pattern := "refresh_token:*"

	var cursor uint64
	var sessions []SessionInfo

	for {
		var keys []string
		var err error
		keys, cursor, err = config.RedisClient.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return nil, err
		}

		// 遍历找到的键
		for _, key := range keys {
			val, err := config.RedisClient.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			// 检查是否属于当前用户
			if val == strconv.FormatUint(uint64(userID), 10) {
				refreshToken := strings.TrimPrefix(key, "refresh_token:")
				// 获取会话信息
				infoKey := "session_info:" + refreshToken
				infoStr, err := config.RedisClient.Get(ctx, infoKey).Result()
				if err != nil {
					if err != redis.Nil {
						continue
					}
					continue
				}

				var sessionInfo SessionInfo
				if err := json.Unmarshal([]byte(infoStr), &sessionInfo); err != nil {
					continue
				}
				sessions = append(sessions, sessionInfo)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return sessions, nil
}

// SaveSessionInfo 保存会话信息
func SaveSessionInfo(c *gin.Context, refreshToken string, userID uint) error {
	deviceInfo := DeviceInfo{
		UserAgent:   c.GetHeader("User-Agent"),
		IP:          c.ClientIP(),
		LastActive:  time.Now(),
		LastRefresh: time.Now(),
	}

	sessionInfo := SessionInfo{
		SessionID:  generateSessionID(refreshToken),
		DeviceInfo: deviceInfo,
		CreatedAt:  time.Now(),
	}

	infoBytes, err := json.Marshal(sessionInfo)
	if err != nil {
		return err
	}

	// 存储会话信息
	infoKey := "session_info:" + refreshToken
	if err := config.RedisClient.Set(c, infoKey, string(infoBytes), 30*24*time.Hour).Err(); err != nil {
		return err
	}

	// 存储会话ID到refreshToken的映射
	sessionKey := "session:" + sessionInfo.SessionID
	return config.RedisClient.Set(c, sessionKey, refreshToken, 30*24*time.Hour).Err()
}

// RevokeSession 吊销指定会话
func RevokeSession(sessionID string) error {
	ctx := context.Background()

	// 获取refreshToken
	sessionKey := "session:" + sessionID
	refreshToken, err := config.RedisClient.Get(ctx, sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // 会话已经不存在
		}
		return err
	}

	// 删除会话ID映射
	if err := config.RedisClient.Del(ctx, sessionKey).Err(); err != nil {
		return err
	}

	// 调用Logout删除相关信息
	return Logout(refreshToken)
}

// UpdateSessionLastActive 更新会话最后活跃时间
func UpdateSessionLastActive(c *gin.Context, refreshToken string) error {
	ctx := context.Background()
	infoKey := "session_info:" + refreshToken

	infoStr, err := config.RedisClient.Get(ctx, infoKey).Result()
	if err != nil {
		return err
	}

	var sessionInfo SessionInfo
	if err := json.Unmarshal([]byte(infoStr), &sessionInfo); err != nil {
		return err
	}

	sessionInfo.DeviceInfo.LastActive = time.Now()

	infoBytes, err := json.Marshal(sessionInfo)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(ctx, infoKey, string(infoBytes), 30*24*time.Hour).Err()
}

// UpdateSessionLastRefresh 更新会话最后刷新时间和设备信息
func UpdateSessionLastRefresh(c *gin.Context, refreshToken string) error {
	ctx := context.Background()
	infoKey := "session_info:" + refreshToken

	infoStr, err := config.RedisClient.Get(ctx, infoKey).Result()
	if err != nil {
		return err
	}

	var sessionInfo SessionInfo
	if err := json.Unmarshal([]byte(infoStr), &sessionInfo); err != nil {
		return err
	}

	sessionInfo.DeviceInfo.LastRefresh = time.Now()
	sessionInfo.DeviceInfo.UserAgent = c.GetHeader("User-Agent")
	sessionInfo.DeviceInfo.IP = c.ClientIP()

	infoBytes, err := json.Marshal(sessionInfo)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(ctx, infoKey, string(infoBytes), 30*24*time.Hour).Err()
}

// Logout 退出登录
func Logout(refreshToken string) error {
	ctx := context.Background()

	// 获取会话信息以获取sessionID
	infoKey := "session_info:" + refreshToken
	infoStr, err := config.RedisClient.Get(ctx, infoKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	if err != redis.Nil {
		var sessionInfo SessionInfo
		if err := json.Unmarshal([]byte(infoStr), &sessionInfo); err != nil {
			return err
		}
		// 删除会话ID映射
		sessionKey := "session:" + sessionInfo.SessionID
		if err := config.RedisClient.Del(ctx, sessionKey).Err(); err != nil {
			return err
		}
	}

	// 删除会话信息
	if err := config.RedisClient.Del(ctx, infoKey).Err(); err != nil {
		return err
	}

	// 删除刷新令牌
	refreshSessionKey := "refresh_token:" + refreshToken
	if err := config.RedisClient.Del(ctx, refreshSessionKey).Err(); err != nil {
		return err
	}

	return nil
}

// LogoutAllDevices 退出所有设备
func LogoutAllDevices(userID uint) error {
	ctx := context.Background()
	pattern := "refresh_token:*"

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = config.RedisClient.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			val, err := config.RedisClient.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			if val == strconv.FormatUint(uint64(userID), 10) {
				refreshToken := strings.TrimPrefix(key, "refresh_token:")
				if err := Logout(refreshToken); err != nil {
					return err
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
