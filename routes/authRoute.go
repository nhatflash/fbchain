package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
	_ "github.com/nhatflash/fbchain/docs"
)


func AuthRoutes(r *gin.Engine, prefix string, ac *controller.AuthController) {
	auth := r.Group(prefix);

	auth.POST("/signin", ac.SignIn)
	auth.POST("/signup/tenant", ac.TenantSignUp)
}