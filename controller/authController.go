package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	appError "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/service"
	"net/http"
)

type AuthController struct {
	Db   *sql.DB
}

func (authController AuthController) Login(c *gin.Context) {
	var signInRequest client.SignInRequest
	if reqErr := c.ShouldBindJSON(&signInRequest); reqErr != nil {
		c.Error(appError.BadRequestError("Validation failed on fields: " + reqErr.Error()))
		return
	}
	signInResponse, logErr := service.HandleSignIn(&signInRequest, authController.Db, c)
	if logErr != nil {
		c.Error(logErr)
		return
	}
	api.SuccessMessage(http.StatusOK, "Login successfully", signInResponse, c)
}



func (authController AuthController) InitializedTenantRegister(c *gin.Context) {
	var registerTenantRequest client.InitializedTenantRegisterRequest

	if reqErr := c.ShouldBindJSON(&registerTenantRequest); reqErr != nil {
		c.Error(appError.BadRequestError("Validation failed on fields: " + reqErr.Error()))
		return
	}

	userTenantResponse, resErr := service.HandleInitializedTenantRegister(&registerTenantRequest, authController.Db, c)

	if resErr != nil {
		c.Error(resErr)
		return
	}
	api.SuccessMessage(201, "Register successfully", userTenantResponse, c)
}