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
	// 将输入的 role 转换为 int，百炼平台无法区分数字和字符串
	roleInt, err := strconv.Atoi(initReq.Role)
	if err != nil {
		c.JSON(http.StatusOK, model.ResponseERR("Invalid role: "+err.Error(), nil))
		return
	}
	wrapperRole := model.Role(roleInt)
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
