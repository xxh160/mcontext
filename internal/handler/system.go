package handler

import (
    "net/http"
    "mcontext/internal/service"
    "mcontext/internal/model"
    "github.com/gin-gonic/gin"
)

func WarmUp(c *gin.Context) {
    err := service.WarmUp()
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.Response{Code: 1, Message: "Failed to warm up"})
        return
    }
    c.JSON(http.StatusOK, model.Response{Code: 0, Message: "WarmUp successful"})
}

func CoolDown(c *gin.Context) {
    err := service.CoolDown()
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.Response{Code: 1, Message: "Failed to cool down"})
        return
    }
    c.JSON(http.StatusOK, model.Response{Code: 0, Message: "CoolDown successful"})
}
