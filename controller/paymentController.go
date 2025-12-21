package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/service"
	"github.com/shopspring/decimal"
)

type PaymentController struct {
	PaymentService 			service.IPaymentService
}


func NewPaymentController(ps service.IPaymentService) *PaymentController {
	return &PaymentController{
		PaymentService: ps,
	}
}

// @Summary Payment URL API
// @Router /payment [get]
func GetPaymentUrl(c *gin.Context) {
	var err error
	var price decimal.Decimal
	price, err = decimal.NewFromString("100000.00")
	if err != nil {
		c.Error(err)
		return
	}
	var url string
	url, err = service.GetVnPayUrl(price)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusOK, "Get successfully", url, c)
} 


// @Summary Pay order with cash API
// @Accept json
// @Produce json
// @Param request body client.PayOrderWithCash true "PayOrderWithCash body"
// @Success 200 {object}
// @Failure 400 {object}
// @Security BearerAuth
// @Router /payment/cash [post]
func (pc *PaymentController) PayOrderWithCash(c *gin.Context) {
	var err error
	var req client.PayOrderWithCashRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err = pc.PaymentService.HandleCashPayment(*req.OrderId, req.Notes); err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusOK, "Payment successfully.", nil, c)
}