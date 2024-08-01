package repo

import (
    "github.com/go-redis/redis/v8"
    "golang.org/x/net/context"
)

var (
    Rdb *redis.Client
    Ctx = context.Background()
)

func InitializeRedis() {
    Rdb = redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Redis 服务地址
    })
}
