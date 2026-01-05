package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
	"github.com/nhatflash/fbchain/middleware"
)


func RestaurantRoutes(r *gin.Engine, prefix string, rc *controller.RestaurantController) {

	restaurant := r.Group(prefix)

	restaurant.POST("/", middleware.JwtRestHandler(), middleware.RoleBasedHandler("TENANT"), rc.CreateRestaurant)

	restaurant.POST("/:restaurantId/item", middleware.JwtRestHandler(), middleware.RoleBasedHandler("TENANT"), rc.AddNewRestaurantItem)

	restaurant.POST("/:restaurantId/table", middleware.JwtRestHandler(), middleware.RoleBasedHandler("TENANT"), rc.AddNewRestaurantTable)

	restaurant.GET("/:restaurantId/table/:tableId/qrCode", middleware.JwtRestHandler(), middleware.RoleBasedHandler("TENANT"), rc.GetTableQrCode)

	restaurant.GET("/:restaurantId/order/:orderId/confirm", middleware.JwtRestHandler(), middleware.RoleBasedHandler("RESTAURANT_STAFF"), rc.ConfirmCashPaymentForRestaurantOrder)
}