package routes

import (
	"database/sql"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"net/http"
)

func MainRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/", Ping)
	AuthRoutes(r, "/api/auth", db)
}


// @Summary Ping
// @Router / [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, api.ApiResponse {
		Status: http.StatusOK,
		Message: "Server alive.",
		Data: nil,
	})
}
