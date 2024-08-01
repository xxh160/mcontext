package middleware

import (
	"mcontext/internal/model"
	"mcontext/internal/state"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	availRWLock sync.RWMutex
)

func GetCheckAvail(c *gin.Context) {
	if c.Request.Method == "POST" && (c.Request.URL.Path == "/warmup" || c.Request.URL.Path == "/cooldown") {
		availRWLock.Lock()
	} else {
		availRWLock.RLock()
	}

	if c.Request.URL.Path == "/warmup" && state.GetAvail() {
		c.JSON(http.StatusOK, model.ResponseERR("System is already warmed up", nil))
		c.Abort()
		return
	}

	if c.Request.URL.Path != "/warmup" && !state.GetAvail() {
		c.JSON(http.StatusOK, model.ResponseERR("System is not warmed up", nil))
		c.Abort()
		return
	}

	c.Next()
}

func PutAvail(c *gin.Context) {
	c.Next()

	if c.Request.Method == "POST" && (c.Request.URL.Path == "/warmup" || c.Request.URL.Path == "/cooldown") {
		availRWLock.Unlock()
	} else {
		availRWLock.RUnlock()
	}
}
