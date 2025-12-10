package api

import (
	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Status		int				`json:"status"`
	Message 	string			`json:"message"`
	Data 		any				`json:"data"`
}

func SuccessMessage(status int, message string, data any, c *gin.Context) {
	c.JSON(status, ApiResponse {
		Status: status,
		Message: message,
		Data: data,
	})
}

