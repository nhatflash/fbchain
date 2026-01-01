package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
)

func TableRoutes(r *gin.Engine, prefix string, rc *controller.RestaurantController) {
	table := r.Group(prefix)
	table.GET("/:tableId", rc.ShowRestaurantItemsViaQRCode)
	table.GET("/:tableId/session/start", rc.StartTableOrderingSession)
	table.GET("/:tableId/session/end", rc.EndTableOrderingSession)
	table.POST("/:tableId/order", rc.CreateRestaurantOrder)
	table.POST("/:tableId/order/:orderId/cash", rc.PayRestaurantOrderWithCash)
}