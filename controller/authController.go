package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/service"
)

type AuthController struct {
	AuthService 		service.IAuthService
}

func NewAuthController(as service.IAuthService) *AuthController {
	return &AuthController{
		AuthService: as,
	}
}

// @Summary Sign in API
// @Accept json
// @Produce json
// @Param request body client.SignInRequest true "SignIn body"
// @Success 200 {object} client.SignInResponse
// @Failure 400 {object} error
// @Router /auth/signin [post]
func (ac *AuthController) SignIn(c *gin.Context) {
	var err error
	var signInRequest client.SignInRequest
	if err = c.ShouldBindJSON(&signInRequest); err != nil {
		c.Error(err)
		return
	}
	var res *client.SignInResponse
	res, err = ac.AuthService.HandleSignIn(&signInRequest)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusOK, "Login successfully", res, c)
}

// @Summary Tenant sign up API
// @Accept json
// @Produce json
// @Param request body client.TenantSignUpRequest true "TenantSignUp body"
// @Success 200 {object} client.TenantResponse
// @Failure 400 {object} error
// @Router /auth/signup/tenant [post]
func (ac *AuthController) TenantSignUp(c *gin.Context) {
	var tenantSignUpReq client.TenantSignUpRequest
	var err error
	if err = c.ShouldBindJSON(&tenantSignUpReq); err != nil {
		c.Error(err)
		return
	}
	var res *client.TenantResponse
	res, err = ac.AuthService.HandleTenantSignUp(&tenantSignUpReq)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, "Tenant signed up successfully.", res, c)
}


// @Summary Verify change password OTP API
// @Accept json
// @Produce json
// @Param request body client.VerifyChangePasswordRequest true "VerifyChangePassword body"
// @Success 200 {object} string 
// @Failure 400 {object} error
// @Security BearerAuth
// @Router /auth/change-password/verify [post]
func (ac *AuthController) VerifyChangePassword(c *gin.Context) {
	var verifyChangePasswordReq client.VerifyChangePasswordRequest
	var err error
	if err = c.ShouldBindJSON(&verifyChangePasswordReq); err != nil {
		c.Error(err)
		return
	}
	var res string
	res, err = ac.AuthService.HandleVerifyChangePassword(&verifyChangePasswordReq, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	if res == "Accepted" {
		api.SuccessMessage(http.StatusOK, "Verify successfully.", res, c)
	} else {
		api.SuccessMessage(http.StatusNotAcceptable, "Verify failed.", res, c)
	}
}

// @Summary Get change password OTP API
// @Success 200 {object} string
// @Failure 400 {object} error
// @Security BearerAuth
// @Router /auth/change-password/verify [get]
func (ac *AuthController) GetChangePasswordVerifiedOTP(c *gin.Context) {
	res, err := ac.AuthService.GenerateChangePasswordVerifyOTP(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusOK, "Verify change password OTP sent successfully.", res, c)
}
