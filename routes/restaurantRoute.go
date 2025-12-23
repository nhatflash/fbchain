package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
	"github.com/nhatflash/fbchain/middleware"
)


func RestaurantRoutes(r *gin.Engine, prefix string, rc *controller.RestaurantController) {

	restaurant := r.Group(prefix, middleware.JwtRestHandler(), middleware.RoleBasedHandler("TENANT"))
	restaurant.POST("/", rc.CreateRestaurant)
	restaurant.POST("/:restaurantId/item", rc.AddNewRestaurantItem)
}