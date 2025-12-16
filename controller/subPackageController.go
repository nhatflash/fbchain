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

type SubPackageController struct {
	Db		*sql.DB
}


// @Summary Create subscription package API
// @Accept json
// @Produce json
// @Param request body client.CreateSubPackageRequest true "CreateSubPackage body"
// @Success 201 {object} client.SubPackageResponse
// @Failure 400 {object} error
// @Security BearerAuth
// @Router /admin/subscription [post]
func (sc SubPackageController) CreateSubPackage(c *gin.Context) {
	var createSubPackageReq client.CreateSubPackageRequest
	var err error

	if err = c.ShouldBindJSON(&createSubPackageReq); err != nil {
		c.Error(err)
		return
	}
	var res *client.SubPackageResponse
	subPackageService := service.NewSubPackageService(sc.Db)
	res, err = subPackageService.HandleCreateSubPackage(&createSubPackageReq)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, "Subscription package created successfully.", res, c)
}