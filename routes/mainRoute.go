package routes

import (
	"net/http"

	"github.com/nhatflash/fbchain/controller"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
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
	r.GET("/", Ping)
	AuthRoutes(r, "/api/auth", ac)
	AdminRoutes(r, "/api/admin", spc)
	TenantRoutes(r, "/api/tenant", tc, rc, oc)
	ProfileRoutes(r, "/api/profile", uc)
	PaymentRoutes(r, "/api/payment", pc)
	RestaurantRoutes(r, "/api/restaurant", rc)
	r.GET("/api/table/:tableId", rc.ShowRestaurantItemsViaQRCode)
}



func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, api.ApiResponse {
		Status: http.StatusOK,
		Message: "Server alive.",
		Data: nil,
	})
}
