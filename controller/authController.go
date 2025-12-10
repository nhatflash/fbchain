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
// @Router /auth/login [post]
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


// @Summary Intial tenant sign up API
// @Accept json
// @Produce json
// @Param request body client.InitialTenantRegisterRequest true "InitialTenantRegister body"
// @Success 201 {object} client.UserResponse
// @Failure 400 {object} error
// @Router /auth/register/tenant/initial [post]
func (authController AuthController) InitialTenantRegister(c *gin.Context) {
	var initialTenantRegisterReq client.InitialTenantRegisterRequest

	if reqErr := c.ShouldBindJSON(&initialTenantRegisterReq); reqErr != nil {
		c.Error(reqErr)
		return
	}

	userTenantRes, resErr := service.HandleInitialTenantRegister(&initialTenantRegisterReq, authController.Db, c)

	if resErr != nil {
		c.Error(resErr)
		return
	}
	api.SuccessMessage(http.StatusCreated, "Initial registration successfully", userTenantRes, c)
}


// @Summary Completed tenant register API
// @Accept json
// @Produce json
// @Param request body client.CompletedTenantRegisterRequest true "CompletedTenantRegister body"
// @Success 200 {object} client.UserResponse
// @Failure 400 {object} error
// @Security BearerAuth
// @Router /auth/register/tenant/completed [post]
func (authController AuthController) CompletedTenantRegister(c *gin.Context) {
	var completedTenantRegisterReq client.CompletedTenantRegisterRequest

	if reqErr := c.ShouldBindJSON(&completedTenantRegisterReq); reqErr != nil {
		c.Error(reqErr)
		return
	}

	userTenantRes, err := service.HandleCompletedTenantRegister(&completedTenantRegisterReq, authController.Db, c)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusOK, "Completed registration successfully", userTenantRes, c)
}