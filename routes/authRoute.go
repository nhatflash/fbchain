package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/middleware"
)


func AuthRoutes(router *gin.Engine, prefix string, db *sql.DB) {
	authController := controller.AuthController{
		Db : db,
	}
	auth := router.Group(prefix);

	auth.POST("/login", authController.SignIn)
	auth.POST("/register/tenant/initial", authController.InitialTenantRegister)
	auth.POST("/register/tenant/completed", middleware.JwtAccessHandler(), middleware.RoleBasedHandler("TENANT"), authController.CompletedTenantRegister)
}