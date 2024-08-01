package handler

import (
	"mcontext/internal/model"
	"mcontext/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	service service.SystemService
}

func (h *SystemHandler) Reset(c *gin.Context) {
	err := h.service.Reset(c)
	if err != nil {
		c.JSON(http.StatusOK, model.ResponseERR("Reset failed: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.ResponseOK(nil))
}

func NewSystemHandler(service service.SystemService) *SystemHandler {
	return &SystemHandler{}
}
