package redis

import "github.com/go-redis/redis"

var client *redis.Client

func InitLike() {
	client = redis.NewClient(&redis.Options{
		Addr:     "43.138.59.103:30260",
		Password: "123456", // Redis数据库没有密码
		DB:       0,        // 默认数据库为0
	})
}
