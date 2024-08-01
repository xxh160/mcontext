package middleware

import (
	"mcontext/internal/erro"
	"mcontext/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorResolve() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先调用 c.Next 执行后面的中间件
		c.Next()

		// 所有中间件及 router 处理完毕后从这里开始执行
		// 检查 c.Errors 中是否有错误
		if len(c.Errors) <= 0 {
			return
		}

		// 处理最后一个 error
		lastErr := c.Errors.Last().Err
		// 若是自定义的错误则将 msg 返回
		// 非自定义错误则返回详细错误信息 err.Error
		if customError, ok := lastErr.(*erro.CustomError); ok {
			c.JSON(http.StatusOK, model.ResponseERR(customError.Msg, nil))
		} else {
			c.JSON(http.StatusOK, model.ResponseERR(lastErr.Error(), nil))
		}
	}
}
