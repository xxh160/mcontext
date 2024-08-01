package repo

import (
	"context"
	"log"
	"mcontext/internal/conf"

	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: conf.RedisAddr,
	})

	// 可连接性检测
	pong, err := rdb.Ping(context.Background()).Result()
	log.Printf("Redis: %s\n", pong)
	return rdb, err
}
