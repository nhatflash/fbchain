package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/service"
)

type AuthController struct {
	Db   *sql.DB
}

// @Summary Sign in API
// @Accept json
// @Produce json
// @Param request body client.SignInRequest true "SignIn" body
// @Success 200 {object} client.SignInResponse
// @Failure 400 {object} error
// @Router /auth/signin [post]
func (authController AuthController) SignIn(c *gin.Context) {
	var signInRequest client.SignInRequest
	if reqErr := c.ShouldBindJSON(&signInRequest); reqErr != nil {
		c.Error(reqErr)
		return
	}
	signInResponse, logErr := service.HandleSignIn(&signInRequest, authController.Db, c)
	if logErr != nil {
		c.Error(logErr)
		return
	}
	api.SuccessMessage(http.StatusOK, "Login successfully", signInResponse, c)
}

// @Summary Tenant sign up API
// @Accept json
// @Produce json
// @Param request body client.TenantSignUpRequest true "TenantSignUp body"
// @Success 200 {object} client.TenantResponse
// @Failure 400 {object} error
// @Router /auth/signup/tenant [post]
func (authController AuthController) TenantSignUp(c *gin.Context) {
	var tenantSignUpReq client.TenantSignUpRequest
	if reqErr := c.ShouldBindJSON(&tenantSignUpReq); reqErr != nil {
		c.Error(reqErr)
		return
	}
	res, resErr := service.HandleTenantSignUp(&tenantSignUpReq, authController.Db)
	if resErr != nil {
		c.Error(resErr)
		return
	}
	api.SuccessMessage(http.StatusCreated, "Tenant signed up successfully.", res, c)
	
}