package controller

import (
	"net/http"
	"github.com/nhatflash/fbchain/constant"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/service"
	_ "github.com/nhatflash/fbchain/docs"
)

type SubPackageController struct {
	SubPackageService 			service.ISubPackageService
}

func NewSubPackageController(sps service.ISubPackageService) *SubPackageController {
	return &SubPackageController{
		SubPackageService: sps,
	}
}


// @Summary Create Subscription Package API
// @Accept json
// @Produce json
// @Param request body client.CreateSubPackageRequest true "CreateSubPackage body"
// @Security BearerAuth
// @Router /admin/subscription [post]
func (spc *SubPackageController) CreateSubPackage(c *gin.Context) {
	var req client.CreateSubPackageRequest
	var err error

	if err = c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	var res *client.SubPackageResponse
	res, err = spc.SubPackageService.HandleCreateSubPackage(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, constant.SUBSCRIPTION_CREATED_SUCCESS, res, c)
}