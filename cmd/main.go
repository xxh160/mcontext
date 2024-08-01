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
	memoryExitFunc func()
	topicExitFunc  func()
	rdb            *redis.Client
)

func serverInit() (*gin.Engine, error) {
	log.Println("Server initing...")

	var err error

	// 配置 Redis
	rdb, err = repo.InitRedis()
	if err != nil {
		return nil, err
	}

	// 配置中间件和 API
	// 初始化各种 repo、service、handler
	var engine *gin.Engine
	engine, memoryExitFunc, topicExitFunc = router.InitRouter(rdb)

	// 标识服务为可用
	util.InitState()

	return engine, nil
}

func serverExit(server *http.Server) {
	// 优雅关闭
	memoryExitFunc()
	topicExitFunc()

	err := rdb.Close()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Releases resources if shutdown completes before timeout elapses
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %s", err)
	}

	log.Println("Server exiting...")
}

func main() {
	engine, err := serverInit()
	if err != nil {
		log.Fatalf("Internal error: %s\n", err)
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %s", err)
		}

		log.Println("Server listening...")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号
	<-quit
	// 接收到信号

	serverExit(&server)
}
