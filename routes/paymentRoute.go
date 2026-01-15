package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/controller"
)


func PaymentRoutes(r *gin.Engine, prefix string, pc *controller.PaymentController) {
	payment := r.Group(prefix)
	payment.GET("/vnpay/:orderId", pc.GetVnPayPaymentUrl)
	payment.POST("/cash", pc.PayOrderWithCash)
}