package api

import (
	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Status		int
	Message 	string
	Data 		any
}

func SuccessMessage(status int, message string, data any, c *gin.Context) {
	c.JSON(status, ApiResponse {
		Status: status,
		Message: message,
		Data: data,
	})
}

