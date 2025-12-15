package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	_ "github.com/nhatflash/fbchain/docs"
)

func MainRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/", Ping)
	AuthRoutes(r, "/api/auth", db)
	AdminRoutes(r, "/api/admin", db)
	TenantRoutes(r, "/api/tenant", db)
}



func Ping(c *gin.Context) {
	fmt.Printf("Client IP: %s\n", c.ClientIP())
	c.JSON(http.StatusOK, api.ApiResponse {
		Status: http.StatusOK,
		Message: "Server alive.",
		Data: nil,
	})
}
