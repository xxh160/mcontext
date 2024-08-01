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

func InitRouter(rdb *redis.Client) (*gin.Engine, func(), func()) {
	r := gin.Default()
	r.Use(middleware.State())
	r.Use(middleware.ErrorResolve())

	// Repo
	topicRepo := repo.NewTopicRepo(rdb)
	memoryRepo := repo.NewMemoryRepo(rdb, topicRepo)

	// Service
	topicService := service.NewTopicservice(topicRepo)
	memoryService := service.NewMemoryService(memoryRepo, topicService)
	systemService := service.NewSystemService(topicService, memoryService)

	// Handler
	memoryHandler := handler.NewMemoryHandler(memoryService)
	systemHanler := handler.NewSystemHandler(systemService)

	r.POST("/reset", systemHanler.Reset)
	r.POST("/memory/create", memoryHandler.CreateMemory)
	r.GET("/memory/:debateTag", memoryHandler.GetMemory)
	r.POST("/memory/:debateTag/update", memoryHandler.UpdateMemory)

	return r, func() {
			ctx := context.Background()
			memoryService.Exit(ctx)
		}, func() {
			ctx := context.Background()
			topicService.RemoveTopicDatas(ctx)
		}
}
