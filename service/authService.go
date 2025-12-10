package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/helper"
	appError "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/security"
	"time"
)

func HandleLogin(loginReq *client.LoginRequest, db *sql.DB, c *gin.Context) (*client.SignInResponse, error) {
	user, userErr := repository.GetSignInUser(loginReq.Login, loginReq.Password, db)
	if userErr != nil {
		return nil, userErr
	}
	accessToken, accessErr := security.GenerateJwtAccessToken(user)
	if accessErr != nil {
		return nil, accessErr
	}
	refreshToken, refreshErr := security.GenerateJwtRefreshToken(user)
	if refreshErr != nil {
		return nil, refreshErr
	}
	res := &client.SignInResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		LastLogin: time.Now(),
	}
	return res, nil
}

func HandleInitializedTenantRegister(registerTenantReq *client.InitializedTenantRegisterRequest, db *sql.DB, c *gin.Context) (*client.UserResponse, error) {
	email := registerTenantReq.Email
	firstName := registerTenantReq.FirstName
	lastName := registerTenantReq.LastName
	password := registerTenantReq.Password
	confirmPassword := registerTenantReq.ConfirmPassword
	gender := registerTenantReq.Gender
	birthdateStr := registerTenantReq.Birthdate

	if repository.CheckUserEmailExists(email, db) {
		return nil, appError.BadRequestError("User with this email already exists.")
	}
	if (password != confirmPassword) {
		return nil, appError.BadRequestError("Confirm password does not match.")
	}
	birthdate, dateErr := helper.ConvertToDate(birthdateStr)
	if dateErr != nil {
		return nil, errors.New(dateErr.Error())
	}
	newTenant, dbErr := repository.RegisterTenant(email, firstName, lastName, password, &gender, birthdate, db)
	if dbErr != nil {
		return nil, dbErr
	}
	return helper.MapToUserResponse(newTenant), nil
}