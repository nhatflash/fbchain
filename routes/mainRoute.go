package routes

import (
	"fmt"
	"net/http"

	"github.com/nhatflash/fbchain/controller"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	_ "github.com/nhatflash/fbchain/docs"
)

func MainRoutes(r *gin.Engine, 
				ac *controller.AuthController, 
				spc *controller.SubPackageController, 
				rc *controller.RestaurantController, 
				oc *controller.OrderController, 
				uc *controller.UserController, 
				pc *controller.PaymentController) {
	r.GET("/", Ping)
	r.GET("/api/payment", controller.GetPaymentUrl)
	AuthRoutes(r, "/api/auth", ac)
	AdminRoutes(r, "/api/admin", spc)
	TenantRoutes(r, "/api/tenant", rc, oc)
	ProfileRoutes(r, "/api/profile", uc)
	PaymentRoutes(r, "/api/payment", pc)
}



func Ping(c *gin.Context) {
	fmt.Printf("Client IP: %s\n", c.ClientIP())
	c.JSON(http.StatusOK, api.ApiResponse {
		Status: http.StatusOK,
		Message: "Server alive.",
		Data: nil,
	})
}
