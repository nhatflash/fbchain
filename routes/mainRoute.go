package routes

import (
	"github.com/nhatflash/fbchain/controller"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	_ "github.com/nhatflash/fbchain/docs"
)

func MainRoutes(r *gin.Engine, 
				ac *controller.AuthController, 
				spc *controller.SubPackageController, 
				rc *controller.RestaurantController, 
				oc *controller.OrderController, 
				uc *controller.UserController) {
	r.GET("/", Ping)
	AuthRoutes(r, "/api/auth", ac)
	AdminRoutes(r, "/api/admin", spc)
	TenantRoutes(r, "/api/tenant", rc, oc)
	ProfileRoutes(r, "/api/profile", uc)
}



func Ping(c *gin.Context) {
	fmt.Printf("Client IP: %s\n", c.ClientIP())
	c.JSON(http.StatusOK, api.ApiResponse {
		Status: http.StatusOK,
		Message: "Server alive.",
		Data: nil,
	})
}
