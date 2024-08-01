package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "mcontext/internal/service"
    "mcontext/internal/model"
)

func InitMemory(c *gin.Context) {
    var initReq model.InitRequest
    if err := c.BindJSON(&initReq); err != nil {
        c.JSON(http.StatusBadRequest, model.Response{Code: 1, Message: "Invalid request body"})
        return
    }

    debateMemory, err := service.InitMemory(initReq.Topic, initReq.Role, initReq.Question)
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.Response{Code: 1, Message: "Failed to initialize memory"})
        return
    }

    c.JSON(http.StatusOK, model.Response{Code: 0, Message: "InitMemory successful", Data: debateMemory})
}

func GetMemory(c *gin.Context) {
    debateTag, err := strconv.Atoi(c.Param("debateTag"))
    if err != nil {
        c.JSON(http.StatusBadRequest, model.Response{Code: 1, Message: "Invalid debateTag"})
        return
    }

    debateMemory, err := service.GetMemory(debateTag)
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.Response{Code: 1, Message: "Failed to get memory"})
        return
    }

    c.JSON(http.StatusOK, model.Response{Code: 0, Message: "GetMemory successful", Data: debateMemory})
}

func UpdateMemory(c *gin.Context) {
    debateTag, err := strconv.Atoi(c.Param("debateTag"))
    if err != nil {
        c.JSON(http.StatusBadRequest, model.Response{Code: 1, Message: "Invalid debateTag"})
        return
    }

    var updateReq model.UpdateRequest
    if err := c.BindJSON(&updateReq); err != nil {
        c.JSON(http.StatusBadRequest, model.Response{Code: 1, Message: "Invalid request body"})
        return
    }

    err = service.UpdateMemory(debateTag, updateReq.Dialog, updateReq.Last)
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.Response{Code: 1, Message: "Failed to update memory"})
        return
    }

    c.JSON(http.StatusOK, model.Response{Code: 0, Message: "UpdateMemory successful"})
}
