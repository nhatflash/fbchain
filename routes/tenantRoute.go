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
	tenant := r.Group(prefix, middleware.JwtRestHandler(), middleware.RoleBasedHandler("TENANT"))
	tenant.POST("/restaurant", middleware.JwtRestHandler(), middleware.RoleBasedHandler("TENANT"), restaurantController.CreateRestaurant)
	tenant.POST("/order", orderController.PaySubPackage)
}
