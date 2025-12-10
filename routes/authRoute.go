package routes

import (
	"database/sql"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
)


func AuthRoutes(router *gin.Engine, prefix string, db *sql.DB) {
	authController := controller.AuthController{
		Db : db,
	}
	auth := router.Group(prefix);

	auth.POST("/login", authController.SignIn)
	auth.POST("/register/tenant/initial", authController.InitialTenantRegister)
}