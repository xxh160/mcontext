package handler

import (
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
	var initReq model.InitRequest
	if err := c.BindJSON(&initReq); err != nil {
		c.JSON(http.StatusOK, model.ResponseERR("Invalid request body: "+err.Error(), nil))
		return
	}

	// 去除输入的辩题的空格
	topic := strings.TrimSpace(initReq.Topic)
	wrapperRole := model.Role(initReq.Role)
	debateMemory, err := h.service.CreateMemory(c, topic, wrapperRole, initReq.Question)
	if err != nil {
		c.JSON(http.StatusOK, model.ResponseERR("Failed to init DebateMemory: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.ResponseOK(debateMemory))
}

func (h *MemoryHandler) GetMemory(c *gin.Context) {
	debateTag, err := strconv.Atoi(c.Param("debateTag"))
	if err != nil {
		c.JSON(http.StatusOK, model.ResponseERR("Invalid debateTag: "+err.Error(), nil))
		return
	}

	debateMemory, err := h.service.GetMemory(c, debateTag)
	if err != nil {
		c.JSON(http.StatusOK, model.ResponseERR("Failed to get DebateMemory: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.ResponseOK(debateMemory))
}

func (h *MemoryHandler) UpdateMemory(c *gin.Context) {
	debateTag, err := strconv.Atoi(c.Param("debateTag"))
	if err != nil {
		c.JSON(http.StatusOK, model.ResponseERR("Invalid debateTag: "+err.Error(), nil))
		return
	}

	var updateReq model.UpdateRequest
	if err := c.BindJSON(&updateReq); err != nil {
		c.JSON(http.StatusOK, model.ResponseERR("Invalid request body: "+err.Error(), nil))
		return
	}

	err = h.service.UpdateMemory(c, debateTag, updateReq.Dialog, updateReq.Last)
	if err != nil {
		c.JSON(http.StatusOK, model.ResponseERR("Failed to update DebateMemory: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.ResponseOK(nil))
}

func NewMemoryHandler(service service.MemoryService) *MemoryHandler {
	return &MemoryHandler{
		service: service,
	}
}
