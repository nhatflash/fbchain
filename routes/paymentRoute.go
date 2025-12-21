package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
	"github.com/nhatflash/fbchain/middleware"
)


func PaymentRoutes(r *gin.Engine, prefix string, pc *controller.PaymentController) {
	payment := r.Group(prefix, middleware.JwtRestHandler())
	payment.POST("/cash", pc.PayOrderWithCash)
}