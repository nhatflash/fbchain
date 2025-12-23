package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/constant"
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
// @Security BearerAuth
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
	api.SuccessMessage(http.StatusOK, constant.PAYMENT_WITH_CASH_SUCCESS, nil, c)
}



// @Summary Pay order with online payment API
// @Accept json
// @Produce json
// @Param method path client.OnlineMethod true "Online method"
// @Param orderId query string true "Order ID"
// @Security BearerAuth
// @Router /payment/online/{method} [post]
func (pc *PaymentController) PayOrderWithOnlinePayment(c *gin.Context) {
	var err error
	method := c.Param("method")
	orderId := c.Query("orderId")

	var oId int64
	oId, err = strconv.ParseInt(orderId, 10, 64)
	if err != nil {
		c.Error(appErr.BadRequestError("Invalid orderId."))
		return
	}
	var res string
	switch method {
		case "VNPAY":
			res, err = pc.VnPayService.GetOrderVnPayUrl(c.Request.Context(), c.ClientIP(), oId)
			if err != nil {
				c.Error(err)
				return
			}
		default:
			c.Error(appErr.BadRequestError("Invalid method"))
			return
	}
	api.SuccessMessage(http.StatusOK, "Url retrieved successfully.", res, c)
}