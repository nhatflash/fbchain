package routes

import (
	"database/sql"
	"github.com/nhatflash/fbchain/middleware"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
)

func AdminRoutes(r *gin.Engine, prefix string, db *sql.DB) {
	subscriptionController := controller.SubscriptionController{
		Db: db,
	}
	admin := r.Group(prefix, middleware.JwtRestHandler(), middleware.RoleBasedHandler("ADMIN"))

	admin.POST("/subscription", subscriptionController.CreateSubscription)
}