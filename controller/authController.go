package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/service"
)

type AuthController struct {
	Db   *sql.DB
}

func (authController AuthController) Login(c *gin.Context) {
	var loginRequest client.LoginRequest

	if reqErr := c.ShouldBindJSON(&loginRequest); reqErr != nil {
		api.ErrorMessage(400, "VALIDATION_ERROR", "JSON binding error: " + reqErr.Error(), c)
		return
	}
	userEmail, logErr := service.HandleLogin(&loginRequest, authController.Db, c)
	if logErr != nil {
		api.ErrorMessage(401, "LOGIN_FAILED", logErr.Error(), c)
		return
	}
	api.SuccessMessage(200, "Login successfully", userEmail, c)
}



func (authController AuthController) RegisterTenant(c *gin.Context) {
	var registerTenantRequest client.RegisterTenantRequest

	if reqErr := c.ShouldBindJSON(&registerTenantRequest); reqErr != nil {
		api.ErrorMessage(400, "VALIDATION_ERROR", "JSON binding error: " + reqErr.Error(), c)
		return
	}

	userTenantResponse, resErr := service.HandleRegisterTenant(&registerTenantRequest, authController.Db, c)

	if resErr != nil {
		api.ErrorMessage(400, "REGISTER_ERROR", resErr.Error(), c)
		return
	}
	api.SuccessMessage(201, "Register successfully", userTenantResponse, c)
}