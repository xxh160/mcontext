package repo

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis 服务地址
	})

	// 可连接性检测
	pong, err := rdb.Ping(context.Background()).Result()
	log.Printf("Pong: %s\n", pong)

	return rdb, err
}
