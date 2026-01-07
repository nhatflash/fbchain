package routes

import (
	"github.com/nhatflash/fbchain/middleware"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
)	

func AdminRoutes(r *gin.Engine, prefix string, spc *controller.SubPackageController, adc *controller.AdminController) {
	admin := r.Group(prefix, middleware.RoleBasedHandler("ADMIN"))

	admin.POST("/subscription", spc.CreateSubPackage)
	admin.POST("/signup/staff", adc.CreateNewStaff)
}