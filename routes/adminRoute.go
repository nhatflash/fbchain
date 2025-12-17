package routes

import (
	"github.com/nhatflash/fbchain/middleware"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
)

func AdminRoutes(r *gin.Engine, prefix string, spc *controller.SubPackageController) {
	admin := r.Group(prefix, middleware.JwtRestHandler(), middleware.RoleBasedHandler("ADMIN"))

	admin.POST("/subscription", spc.CreateSubPackage)
}