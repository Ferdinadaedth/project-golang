package redis

import "github.com/go-redis/redis"

var client *redis.Client

func InitLike() {
	client = redis.NewClient(&redis.Options{
		Addr:     "47.108.208.111:6379",
		Password: "h74o+JIi5SpSY3MU", // Redis数据库没有密码
		DB:       0,                  // 默认数据库为0
	})
}
