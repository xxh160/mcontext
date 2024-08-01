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
)

var (
	memoryExitFunc func()
	topicExitFunc  func()
)

func serverInit() *gin.Engine {
	// 配置 Redis
	rdb := repo.InitRedis()

	// 配置中间件和 API
	// 初始化各种 repo、service、handler
	var engine *gin.Engine
	engine, memoryExitFunc, topicExitFunc = router.InitRouter(rdb)

	// 标识服务为可用
	util.InitState()

	return engine
}

func serverExit(server *http.Server) {
	// 优雅关闭
	memoryExitFunc()
	topicExitFunc()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Releases resources if shutdown completes before timeout elapses
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %s", err)
	}

	log.Println("Server exiting...")
}

func main() {
	engine := serverInit()

	server := http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号
	<-quit
	// 接收到信号

	serverExit(&server)
}
