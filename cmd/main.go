package main

import (
	"github.com/gin-gonic/gin"
	"mcontext/internal/handler"
	"mcontext/internal/middleware"
	"mcontext/internal/repo"
	_ "mcontext/internal/state"
)

var Avail bool

func main() {
	repo.InitializeRedis()

	r := gin.Default()

	r.Use(middleware.GetCheckAvail)
	r.Use(middleware.PutAvail)

	r.POST("/warmup", handler.WarmUp)
	r.POST("/cooldown", handler.CoolDown)
	r.POST("/memory/init", handler.InitMemory)
	r.GET("/memory/:debateTag", handler.GetMemory)
	r.POST("/memory/:debateTag/update", handler.UpdateMemory)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
