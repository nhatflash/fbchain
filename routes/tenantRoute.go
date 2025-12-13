package routes

import (
	"database/sql"
	"github.com/nhatflash/fbchain/controller"
	"github.com/nhatflash/fbchain/middleware"
	"github.com/gin-gonic/gin"
)

func TenantRoutes(r *gin.Engine, prefix string, db *sql.DB) {
	restaurantController := controller.RestaurantController{
		Db: db,
	}
	tenant := r.Group(prefix, middleware.JwtAccessHandler(), middleware.RoleBasedHandler("TENANT"))
	tenant.POST("/restaurant", restaurantController.CreateRestaurant)
}