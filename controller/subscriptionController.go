package controller

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/service"
	_ "github.com/nhatflash/fbchain/docs"
)

type SubscriptionController struct {
	Db		*sql.DB
}


// @Summary Create subscription API
// @Accept json
// @Produce json
// @Param request body client.CreateSubscriptionRequest true "CreateSubscription body"
// @Success 201 {object} client.SubscriptionResponse
// @Failure 400 {object} error
// @Security BearerAuth
// @Router /admin/subscription [post]
func (sc SubscriptionController) CreateSubscription(c *gin.Context) {
	var createSubscriptionReq client.CreateSubscriptionRequest
	var err error

	if err = c.ShouldBindJSON(&createSubscriptionReq); err != nil {
		c.Error(err)
		return
	}
	var res *client.SubscriptionResponse
	subscriptionService := service.NewSubscriptionService(sc.Db)
	res, err = subscriptionService.HandleCreateSubscription(&createSubscriptionReq)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, "Subscription created successfully.", res, c)
}