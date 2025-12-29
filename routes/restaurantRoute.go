package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
)


func RestaurantRoutes(r *gin.Engine, prefix string, rc *controller.RestaurantController) {

	restaurant := r.Group(prefix)
	restaurant.POST("/", rc.CreateRestaurant)
	restaurant.POST("/:restaurantId/item", rc.AddNewRestaurantItem)
	restaurant.POST("/:restaurantId/table", rc.AddNewRestaurantTable)
	restaurant.GET("/test", rc.GetTableQrCode)
}