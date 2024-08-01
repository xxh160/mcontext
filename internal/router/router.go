package router

import (
	"context"
	"mcontext/internal/repo"
	"mcontext/internal/router/handler"
	"mcontext/internal/router/middleware"
	"mcontext/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ExitFunc func() error
type InitFunc func() error

func InitRouter(rdb *redis.Client) (*gin.Engine, InitFunc, InitFunc, ExitFunc, ExitFunc) {
	r := gin.Default()

	r.Use(middleware.State())
	r.Use(middleware.ErrorResolve())

	// Repo
	topicRepo := repo.NewTopicRepo(rdb)
	memoryRepo := repo.NewMemoryRepo(rdb, topicRepo)

	// Service
	topicService := service.NewTopicService(topicRepo)
	memoryService := service.NewMemoryService(memoryRepo, topicService)
	systemService := service.NewSystemService(topicService, memoryService)

	// Handler
	memoryHandler := handler.NewMemoryHandler(memoryService)
	systemHanler := handler.NewSystemHandler(systemService)

	r.POST("/reset", systemHanler.Reset)
	r.POST("/memory/create", memoryHandler.CreateMemory)
	r.GET("/memory", memoryHandler.GetMemory)
	r.POST("/memory/update", memoryHandler.UpdateMemory)

	memoryInit := func() error {
		ctx := context.Background()
		return memoryService.Init(ctx)
	}

	memoryExit := func() error {
		ctx := context.Background()
		return memoryService.Exit(ctx)
	}

	topicInit := func() error {
		ctx := context.Background()
		return topicService.LoadTopicDatas(ctx)
	}

	topicExit := func() error {
		ctx := context.Background()
		return topicService.RemoveTopicDatas(ctx)
	}

	return r, memoryInit, topicInit, memoryExit, topicExit
}
