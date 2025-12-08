package routes

import (
	"github.com/gin-gonic/gin"
	api "github.com/nhatflash/fbchain/api"
)

func MainRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, api.ApiResponse {
			Code: 200,
			Message: "Server alive.",
			Data: nil,
		})
	})
}