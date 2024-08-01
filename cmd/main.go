package main

import (
	"context"
	"log"
	"mcontext/internal/repo"
	"mcontext/internal/router"
	"mcontext/internal/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	memoryExit func() error
	topicExit  func() error
	rdb        *redis.Client
)

func serverInit() (*gin.Engine, error) {
	log.Printf("Server initing...\n")

	var err error

	// 配置 Redis
	rdb, err = repo.InitRedis()
	if err != nil {
		log.Printf("Cannot init redis: %s\n", err)
		return nil, err
	}

	// 配置中间件和 API
	// 初始化各种 repo、service、handler
	var engine *gin.Engine
	var memoryInit func() error
	var topicInit func() error
	engine, memoryInit, topicInit, memoryExit, topicExit = router.InitRouter(rdb)

	// 初始化 memory service
	err = memoryInit()
	if err != nil {
		log.Printf("Cannot init memory service: %s\n", err)
		return nil, err
	}

	// 初始化 topic service
	err = topicInit()
	if err != nil {
		log.Printf("Cannot init topic service: %s\n", err)
		memoryExit()
		return nil, err
	}

	// 标识服务为可用
	util.InitState()

	return engine, nil
}

func serverExit(server *http.Server) {
	// 关闭 memory service
	if err := memoryExit(); err != nil {
		log.Printf("Cannot exit memory service: %s\n", err)
	}

	// 关闭 topic service
	if err := topicExit(); err != nil {
		log.Printf("Cannot exit topic service: %s\n", err)
	}

	err := rdb.Close()
	if err != nil {
		log.Printf("Cannot close redis: %s\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Releases resources if shutdown completes before timeout elapses
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %s\n", err)
	}

	log.Printf("Server exiting...\n")
}

func main() {
	engine, err := serverInit()
	if err != nil {
		log.Printf("Server init error: %s\n", err)
		return
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server listen error: %s\n", err)
		}
		log.Printf("Server listening...\n")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号
	<-quit
	// 接收到信号

	serverExit(&server)
}
