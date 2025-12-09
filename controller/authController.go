package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/service"
	"net/http"
)

type AuthController struct {
	Db   *sql.DB
}

func (authController AuthController) Login(c *gin.Context) {
	var loginRequest client.LoginRequest

	if reqErr := c.ShouldBindJSON(&loginRequest); reqErr != nil {
		c.Error(reqErr)
		return
	}
	userEmail, logErr := service.HandleLogin(&loginRequest, authController.Db, c)
	if logErr != nil {
		c.Error(logErr)
		return
	}
	api.SuccessMessage(http.StatusOK, "Login successfully", userEmail, c)
}



func (authController AuthController) RegisterTenant(c *gin.Context) {
	var registerTenantRequest client.RegisterTenantRequest

	if reqErr := c.ShouldBindJSON(&registerTenantRequest); reqErr != nil {
		c.Error(reqErr)
		return
	}

	userTenantResponse, resErr := service.HandleRegisterTenant(&registerTenantRequest, authController.Db, c)

	if resErr != nil {
		c.Error(resErr)
		return
	}
	api.SuccessMessage(201, "Register successfully", userTenantResponse, c)
}