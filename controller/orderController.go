package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/service"
)

type OrderController struct {
	Db *sql.DB
}

// @Summary Pay subscription API
// @Accept json
// @Produce json
// @Param request body client.PaySubscriptionRequest true "PaySubscription body"
// @Success 201 {object} client.OrderResponse
// @Failure 400 {object} error
// @Security BearerAuth
// @Router /tenant/order [post]
func (oc OrderController) PaySubscription(c *gin.Context) {
	var paySubscriptionReq client.PaySubscriptionRequest
	var err error

	if err = c.ShouldBindJSON(&paySubscriptionReq); err != nil {
		c.Error(err)
		return
	}

	var currTenant *model.Tenant
	currTenant, err = service.GetCurrentTenant(c, oc.Db)
	if err != nil {
		c.Error(err)
		return
	}

	var res *client.OrderResponse
	res, err = service.HandlePaySubscription(&paySubscriptionReq, currTenant.Id, oc.Db)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, "Order created successfully.", res, c)
}
