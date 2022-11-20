package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json: "message"`
}

type statusResponse struct {
	Status string `json: "status"`
}

// в случае ошибки пишем логи и предотвращаем вызов ожидающих обработчиков 
func newErrorResponse(c *gin.Context, statuscode int, massage string) {
	logrus.Error(massage)
	c.AbortWithStatusJSON(statuscode, errorResponse{massage})
}