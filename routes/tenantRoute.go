package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
	"github.com/nhatflash/fbchain/middleware"
)

func TenantRoutes(r *gin.Engine, prefix string, tc *controller.TenantController, rc *controller.RestaurantController, oc *controller.OrderController) {
	tenant := r.Group(prefix, middleware.JwtRestHandler(), middleware.RoleBasedHandler("TENANT"))
	tenant.POST("/verify", tc.CompleteTenantInfo)
	tenant.POST("/order", oc.PaySubPackage)
}
