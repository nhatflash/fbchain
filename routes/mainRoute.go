package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	api "github.com/nhatflash/fbchain/api"
)

func MainRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, api.ApiResponse {
			Code: 200,
			Message: "Server alive.",
			Data: nil,
		})
	})
	AuthRoutes(router, "/api/auth", db)
}
