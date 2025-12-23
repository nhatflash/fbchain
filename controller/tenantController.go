package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/service"
	appErr "github.com/nhatflash/fbchain/error"
)

type TenantController struct {
	TenantService 			service.ITenantService
	UserService 			service.IUserService
}

func NewTenantController(ts service.ITenantService, us service.IUserService) *TenantController {
	return &TenantController{
		TenantService: ts,
		UserService: us,
	}
}


// @Summary Complete Tenant Info API
// @Accept json
// @Produce json
// @Param request body client.TenantInfoRequest true "TenantInfo body"
// @Success 200 {object} client.TenantResponse
// @Failure 400 {object} error
// @Security BearerAuth
// @Router /tenant/verify [post]
func (tc *TenantController) CompleteTenantInfo(c *gin.Context) {
	var req client.TenantInfoRequest
	var err error
	ctx := c.Request.Context()
	if err = c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	var u *model.User
	u, err = tc.UserService.GetCurrentUser(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	if u.IsVerified {
		c.Error(appErr.BadRequestError("Your account is already verified."))
		return
	}

	var res *client.TenantResponse
	res, err = tc.TenantService.HandleCompleteTenantInfo(ctx, u.Id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusOK, "Tenant info completed successfully.", res, c)
}