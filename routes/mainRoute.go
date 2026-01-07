package routes

import (
	"github.com/nhatflash/fbchain/controller"

	"github.com/gin-gonic/gin"
	_ "github.com/nhatflash/fbchain/docs"
)

func MainRoutes(r *gin.Engine, 
				ac *controller.AuthController,
				tc *controller.TenantController, 
				spc *controller.SubPackageController, 
				rc *controller.RestaurantController, 
				oc *controller.OrderController, 
				uc *controller.UserController, 
				pc *controller.PaymentController, 
				adc *controller.AdminController) {
	AuthRoutes(r, "/api/auth", ac)
	AdminRoutes(r, "/api/admin", spc, adc)
	TenantRoutes(r, "/api/tenant", tc, rc, oc)
	ProfileRoutes(r, "/api/profile", uc)
	PaymentRoutes(r, "/api/payment", pc)
	RestaurantRoutes(r, "/api/restaurant", rc, tc)
	TableRoutes(r, "/api/table", rc)
}

