package engine

import (
	"context"
	"fmt"
	"log"
	"mcontext/internal/conf"
	"mcontext/internal/engine/handler"
	"mcontext/internal/engine/middleware"
	"mcontext/internal/repo"
	"mcontext/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	*http.Server
	prepareFunc func() error
	exitFunc    func()
}

func (r *Engine) Prepare() error {
	return r.prepareFunc()
}

func (r *Engine) Exit() {
	r.exitFunc()
}

func New() (*Engine, error) {
	// 配置 Redis
	rdb, err := repo.InitRedis()
	if err != nil {
		return nil, fmt.Errorf("cannot init redis: %w", err)
	}

	// 配置 gin engine
	engine := gin.Default()

	// 配置中间件
	engine.Use(middleware.ErrorResolve())

	// Repo
	topicRepo := repo.NewTopicRepo(rdb)
	memoryRepo := repo.NewMemoryRepo(rdb, topicRepo)

	// Service
	topicService := service.NewTopicService(topicRepo)
	memoryService := service.NewMemoryService(memoryRepo, topicService)

	// Handler
	memoryHandler := handler.NewMemoryHandler(memoryService)

	// 配置 gin 路由
	engine.POST("/memory/create", memoryHandler.CreateMemory)
	engine.GET("/memory", memoryHandler.GetMemory)
	engine.POST("/memory/update", memoryHandler.UpdateMemory)

	// 配置 http 服务器
	r := &Engine{
		Server: &http.Server{
			Addr:    conf.ServerAddr,
			Handler: engine,
		},
	}

	r.prepareFunc = func() error {
		ctx := context.Background()
		err := memoryService.Init(ctx)
		if err != nil {
			return fmt.Errorf("memoryService init: %w", err)
		}

		err = topicService.LoadAllTopicData(ctx)
		if err != nil {
			return fmt.Errorf("topicService init failed: %w", err)
		}

		return nil
	}

	r.exitFunc = func() {
		ctx := context.Background()

		errTopic := topicService.RemoveAllTopicData(ctx)
		if errTopic != nil {
			log.Printf("Cannot remove all topicData: %v", errTopic)
		}

		errMemory := memoryService.Exit(ctx)
		if errMemory != nil {
			log.Printf("Cannot exit memoryService: %v", errMemory)
		}

		errRdb := rdb.Close()
		if errRdb != nil {
			log.Printf("Cannot close redis: %v", errRdb)
		}
	}

	return r, nil
}
