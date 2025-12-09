package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
)

func AuthRoutes(router *gin.Engine, prefix string, db *sql.DB) {
	authController := controller.AuthController{
		Db : db,
	}
	auth := router.Group(prefix, nil);
	auth.POST("/login", authController.Login)
	auth.POST("/register", authController.RegisterTenant)
}