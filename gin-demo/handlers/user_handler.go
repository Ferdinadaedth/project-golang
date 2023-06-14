package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"golandprojects/dao"
	"time"
)

// UserHandler 处理与用户相关的 API 请求
type UserHandler struct {
	RedisClient *redis.Client
}

// NewUserHandler 创建一个 UserHandler 实例
func NewUserHandler(redisClient *redis.Client) *UserHandler {
	return &UserHandler{RedisClient: redisClient}
}

// GetUserByUsername 根据用户名获取用户信息
func (handler *UserHandler) GetUserByUsername(c *gin.Context) {
	// 获取请求参数
	username := c.Param("username")

	// 先尝试从 Redis 中获取数据
	user, err := handler.getFromCache(username)
	if err == nil {
		// 如果缓存中存在，直接返回
		c.JSON(200, gin.H{"data": user})
		return
	}

	// 如果 Redis 中不存在则从数据库获取，然后存储到 Redis 中
	user, err = dao.GetUserByUsername(username)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	// 存储到 Redis 缓存中，设置过期时间为一分钟
	err = handler.saveToCache(username, user, time.Minute)
	if err != nil {
		fmt.Printf("Failed to save user %s to cache: %v\n", username, err)
	}

	c.JSON(200, gin.H{"data": user})
}

// saveToCache 将用户信息存储到 Redis 缓存中
func (handler *UserHandler) saveToCache(username string, user *dao.User, expiration time.Duration) error {
	// 将结构体序列化成 json 字符串
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Failed to marshal user JSON: %v", err)
	}

	// 在 Redis 中保存 json 字符串，并设置过期时间
	res := handler.RedisClient.Set(ctx, fmt.Sprintf("user:%s", username), userJSON, expiration)
	if res.Err() != nil {
		return fmt.Errorf("Failed to save user %s to cache: %v", username, res.Err())
	}

	return nil
}

// getFromCache 从 Redis 缓存中获取用户信息
func (handler *UserHandler) getFromCache(username string) (*dao.User, error) {
	// 在 Redis 中查找对应的值
	res := handler.RedisClient.Get(ctx, fmt.Sprintf("user:%s", username))
	if res.Err() != nil {
		return nil, fmt.Errorf("Failed to get user %s from cache: %v", username, res.Err())
	}

	// 将 json 字符串反序列化成结构体
	var user dao.User
	err := json.Unmarshal([]byte(res.Val()), &user)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal user JSON: %v", err)
	}

	return &user, nil
}
