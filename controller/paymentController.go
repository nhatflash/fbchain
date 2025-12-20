package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/service"
	"github.com/shopspring/decimal"
	_ "github.com/nhatflash/fbchain/docs"
)

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