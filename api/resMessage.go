package api

import (
	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Code		int
	Message 	string
	Data 		any
}

func SuccessMessage(code int, message string, data any, c *gin.Context) {
	c.JSON(code, ApiResponse {
		Code: code,
		Message: message,
		Data: data,
	})
}

