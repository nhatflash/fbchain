package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
	"github.com/nhatflash/fbchain/middleware"
)

func ProfileRoutes(r *gin.Engine, prefix string, uc *controller.UserController) {
	profile := r.Group(prefix, middleware.JwtRestHandler())

	profile.PATCH("/", uc.ChangeUserProfile)
}