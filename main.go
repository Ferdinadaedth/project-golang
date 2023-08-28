package main

import (
	"golandprojects/api"
	"golandprojects/cache"
	"golandprojects/dao/redis"
)

func main() {
	redis.InitLike()
	cache.InitCache()
	api.InitRouter()
}
