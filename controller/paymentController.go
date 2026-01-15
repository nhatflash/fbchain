package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/service"
	appErr "github.com/nhatflash/fbchain/error"
)

type PaymentController struct {
	PaymentService 			service.IPaymentService
	VnPayService 			service.IVnPayService
}


func NewPaymentController(ps service.IPaymentService, vs service.IVnPayService) *PaymentController {
	return &PaymentController{
		PaymentService: ps,
		VnPayService: vs,
	}
}


// @Summary Pay order with cash API
// @Accept json
// @Produce json
// @Param request body client.PayOrderWithCashRequest true "PayOrderWithCash body"
// @Success 200 {object} string
// @Failure 400 {object} error
// @Router /payment/cash [post]
func (pc *PaymentController) PayOrderWithCash(c *gin.Context) {
	var err error
	var req client.PayOrderWithCashRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err = pc.PaymentService.HandleCashPayment(c.Request.Context(), *req.OrderId, req.Notes); err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusOK, "Pay order with cash successfully.", nil, c)
}


// @Summary Get VnPay Payment Url API
// @Param orderId path string true "Order ID"
// @Router /payment/vnpay/{orderId} [get]
func (pc *PaymentController) GetVnPayPaymentUrl(c *gin.Context) {
	var err error
	orderIdParam := c.Param("orderId")
	
	var orderId int64
	orderId, err = strconv.ParseInt(orderIdParam, 10, 64);
	if err != nil {
		c.Error(appErr.BadRequestError("Invalid orderId format."))
		return
	}

	var url string
	clientIp := c.ClientIP()
	url, err = pc.VnPayService.GetOrderVnPayUrl(c.Request.Context(), clientIp, orderId)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusOK, "VnPay url retrieved successfully.", url, c)
}