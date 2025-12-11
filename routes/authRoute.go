package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
	_ "github.com/nhatflash/fbchain/docs"
)


func AuthRoutes(router *gin.Engine, prefix string, db *sql.DB) {
	authController := controller.AuthController{
		Db : db,
	}
	auth := router.Group(prefix);

	auth.POST("/signin", authController.SignIn)
	auth.POST("/signup/tenant", authController.TenantSignUp)
}