package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/middleware"
)


func AuthRoutes(r *gin.Engine, prefix string, ac *controller.AuthController) {
	auth := r.Group(prefix);

	auth.POST("/signin", ac.SignIn)
	auth.POST("/signup", ac.TenantSignUp)
	auth.GET("/change-password/verify", middleware.JwtRestHandler(), ac.GetChangePasswordVerifiedOTP)
	auth.POST("/change-password/verify", middleware.JwtRestHandler(), ac.VerifyChangePassword)
	auth.POST("/change-password", middleware.JwtRestHandler(), ac.ChangePassword)
}