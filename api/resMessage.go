package api

import (
	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Code		int
	Message 	string
	Data 		any
}

type ErrorResponse struct {
	Code		int
	Error		error
	Message		string
}

func SuccessMessage(code int, message string, data any, c *gin.Context) {
	c.JSON(code, ApiResponse {
		Code: code,
		Message: message,
		Data: data,
	})
}

func ErrorMessage(code int, e error, message string, c *gin.Context) {
	c.JSON(code, ErrorResponse {
		Code: code,
		Error: e,
		Message: message,
	})
}

