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
				pc *controller.PaymentController) {
	AuthRoutes(r, "/api/auth", ac)
	AdminRoutes(r, "/api/admin", spc)
	TenantRoutes(r, "/api/tenant", tc, rc, oc)
	ProfileRoutes(r, "/api/profile", uc)
	PaymentRoutes(r, "/api/payment", pc)
	RestaurantRoutes(r, "/api/restaurant", rc)
	r.GET("/api/table/:tableId", rc.ShowRestaurantItemsViaQRCode)
	r.GET("/api/table/:tableId/session/start", rc.StartTableOrderingSession)
	r.GET("/api/table/:tableId/session/end", rc.EndTableOrderingSession)
	r.POST("/api/table/:tableId/order", rc.CreateRestaurantOrder)
}

