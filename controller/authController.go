package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	appError "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/service"
	"net/http"
	_ "github.com/nhatflash/fbchain/docs"
)

type AuthController struct {
	Db   *sql.DB
}

// @Summary Sign in API
// @Accept json
// @Produce json
// @Param request body client.SignInRequest true "SignIn" body
// @Success 200 {object} client.SignInResponse
// @Failure 400 {object} appError.ErrorResponse
// @Router /auth/login [post]
func (authController AuthController) SignIn(c *gin.Context) {
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


// @Summary Intial tenant sign up API
// @Accept json
// @Produce json
// @Param request body client.InitialTenantRegisterRequest true "InitialTenantRegister body"
// @Success 201 {object} client.UserResponse
// @Failure 400 {object} appError.ErrorResponse
// @Router /auth/register/tenant/initial [post]
func (authController AuthController) InitialTenantRegister(c *gin.Context) {
	var registerTenantRequest client.InitialTenantRegisterRequest

	if reqErr := c.ShouldBindJSON(&registerTenantRequest); reqErr != nil {
		c.Error(appError.BadRequestError("Validation failed on fields: " + reqErr.Error()))
		return
	}

	userTenantResponse, resErr := service.HandleInitialTenantRegister(&registerTenantRequest, authController.Db, c)

	if resErr != nil {
		c.Error(resErr)
		return
	}
	api.SuccessMessage(201, "Register successfully", userTenantResponse, c)
}