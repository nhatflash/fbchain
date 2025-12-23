package controller 

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/service"
	"github.com/nhatflash/fbchain/api"
	"net/http"
	_ "github.com/nhatflash/fbchain/docs"
)

type TenantController struct {
	TenantService 			service.ITenantService
}

func NewTenantController(ts service.ITenantService) *TenantController {
	return &TenantController{
		TenantService: ts,
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
	if err = c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	var res *client.TenantResponse
	res, err = tc.TenantService.HandleCompleteTenantInfo(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusOK, "Tenant info completed successfully.", res, c)
}