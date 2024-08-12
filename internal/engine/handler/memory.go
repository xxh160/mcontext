package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"mcontext/internal/model"
	"mcontext/internal/service"

	"github.com/gin-gonic/gin"
)

type MemoryHandler struct {
	service service.MemoryService
}

func (h *MemoryHandler) CreateMemory(c *gin.Context) {
	var createReq model.CreateMemoryRequest
	if err := c.BindJSON(&createReq); err != nil {
		_ = c.Error(fmt.Errorf("invalid request body: %w", err))
		return
	}

	// 去除输入的辩题的空格
	topic := strings.TrimSpace(createReq.Topic)
	debateMemory, err := h.service.CreateMemory(c, topic, createReq.Role, createReq.Question)
	if err != nil {
		_ = c.Error(fmt.Errorf("failed to init DebateMemory: %w", err))
		return
	}

	c.JSON(http.StatusOK, model.ResponseOK(debateMemory))
}

func (h *MemoryHandler) GetMemory(c *gin.Context) {
	debateTag := c.DefaultQuery("debateTag", "-1")
	if debateTag == "-1" {
		_ = c.Error(fmt.Errorf("no debateTag param"))
		return
	}

	// 将 debateTag 转为 int
	tag, err := strconv.Atoi(debateTag)
	if err != nil {
		_ = c.Error(fmt.Errorf("invalid debateTag: %w", err))
		return
	}

	debateMemory, err := h.service.GetMemory(c, tag)
	if err != nil {
		_ = c.Error(fmt.Errorf("failed to get DebateMemory: %w", err))
		return
	}

	c.JSON(http.StatusOK, model.ResponseOK(debateMemory))
}

func (h *MemoryHandler) UpdateMemory(c *gin.Context) {
	var updateReq model.UpdateMemoryRequest
	if err := c.BindJSON(&updateReq); err != nil {
		_ = c.Error(fmt.Errorf("invalid request body: %w", err))
		return
	}

	dialog := model.Dialog{Question: updateReq.Question, Answer: updateReq.Answer}
	last, err := strconv.ParseBool(updateReq.Last)
	if err != nil {
		_ = c.Error(fmt.Errorf("invalid bool last: %w", err))
		return
	}

	err = h.service.UpdateMemory(c, updateReq.DebateTag, dialog, last)
	if err != nil {
		_ = c.Error(fmt.Errorf("failed to update DebateMemory: %w", err))
		return
	}

	c.JSON(http.StatusOK, model.ResponseOK(nil))
}

func NewMemoryHandler(service service.MemoryService) *MemoryHandler {
	return &MemoryHandler{
		service: service,
	}
}
