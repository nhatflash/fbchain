package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
	"github.com/nhatflash/fbchain/middleware"
)

func TenantRoutes(r *gin.Engine, prefix string, db *sql.DB) {
	restaurantController := controller.RestaurantController{
		Db: db,
	}
	orderController := controller.OrderController{
		Db: db,
	}
	tenant := r.Group(prefix, middleware.JwtAccessHandler(), middleware.RoleBasedHandler("TENANT"))
	tenant.POST("/restaurant", restaurantController.CreateRestaurant)
	tenant.POST("/order", orderController.PaySubscription)
}
