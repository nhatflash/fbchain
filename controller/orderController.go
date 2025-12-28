package controller

import (
	"net/http"
	"github.com/nhatflash/fbchain/constant"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/service"
	appErr "github.com/nhatflash/fbchain/error"
)

type OrderController struct {
	OrderService 		service.IOrderService
	TenantService 		service.ITenantService
	UserService 		service.IUserService
}


func NewOrderController(os service.IOrderService, ts service.ITenantService, us service.IUserService) *OrderController {
	return &OrderController{
		OrderService: os,
		TenantService: ts,
		UserService: us,
	}
}

// @Summary Pay subscription package API
// @Accept json
// @Produce json
// @Param request body client.PaySubPackageRequest true "PaySubPackage body"
// @Success 201 {object} client.OrderResponse
// @Failure 400 {object} error
// @Security BearerAuth
// @Router /tenant/order [post]
func (oc OrderController) PaySubPackage(c *gin.Context) {
	var paySubPackageReq client.PaySubPackageRequest
	
	if err := c.ShouldBindJSON(&paySubPackageReq); err != nil {
		c.Error(err)
		return
	}
	ctx := c.Request.Context()
	currUser, err := oc.UserService.FindCurrentUser(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	if !currUser.IsVerified {
		c.Error(appErr.UnauthorizedError("Please verify your account before doing this action."))
		return
	}
	currTenant, err := oc.TenantService.FindTenantByUserId(ctx, currUser.Id)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := oc.OrderService.HandlePaySubPackage(ctx, &paySubPackageReq, currTenant.Id)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, constant.ORDER_CREATED_SUCCESS, res, c)
}
