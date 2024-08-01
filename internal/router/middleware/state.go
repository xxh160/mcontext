package middleware

import (
	"mcontext/internal/model"
	"mcontext/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func State() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" && (c.Request.URL.Path == "/reset") {
			util.ChangeServerStart()
		} else {
			if !util.UseServerStart() {
				c.JSON(http.StatusOK, model.ResponseERR("Server Unavailable", nil))
				c.Abort()
			}
		}

		c.Next()

		if c.Request.Method == "POST" && (c.Request.URL.Path == "/reset") {
			util.ChangeServerEnd()
		} else {
			util.UseServerEnd()
		}

	}
}
