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
	"strings"
	"github.com/nhatflash/fbchain/model"
)


// Sign in method
func HandleSignIn(signInReq *client.SignInRequest, db *sql.DB, c *gin.Context) (*client.SignInResponse, error) {
	login := signInReq.Login

	var loggedUser *model.User
	if strings.Contains(login, "@") {
		user, userErr := repository.GetUserByEmail(signInReq.Login, db)
		if userErr != nil {
			return nil, appError.UnauthorizedError("Invalid credentials.")
		}
		loggedUser = user
	} else {
		user, userErr := repository.GetUserByPhone(signInReq.Login, db)
		if userErr != nil {
			return nil, appError.UnauthorizedError("Invalid credentials.")
		}
		loggedUser = user
	}

	if !security.VerifyPassword(signInReq.Password, loggedUser.Password) {
		return nil, appError.UnauthorizedError("Invalid credentials.")
	}
	
	accessToken, accessErr := security.GenerateJwtAccessToken(loggedUser)
	if accessErr != nil {
		return nil, accessErr
	}
	refreshToken, refreshErr := security.GenerateJwtRefreshToken(loggedUser)
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



// intialized sign up method for tenant
func HandleInitialTenantRegister(initialTenantRegisterReq *client.InitialTenantRegisterRequest, db *sql.DB, c *gin.Context) (*client.UserResponse, error) {
	email := initialTenantRegisterReq.Email
	firstName := initialTenantRegisterReq.FirstName
	lastName := initialTenantRegisterReq.LastName
	password := initialTenantRegisterReq.Password
	confirmPassword := initialTenantRegisterReq.ConfirmPassword
	gender := initialTenantRegisterReq.Gender
	birthdateStr := initialTenantRegisterReq.Birthdate

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
	hashedPassword, hashErr := security.GenerateHashedPassword(password)
	if hashErr != nil {
		return nil, appError.InternalError("An unexpected error when creating hash password")
	}
	newTenant, dbErr := repository.RegisterTenant(email, firstName, lastName, hashedPassword, &gender, birthdate, db)
	if dbErr != nil {
		return nil, dbErr
	}
	return helper.MapToUserResponse(newTenant), nil
}


func HandleCompletedTenantRegister(completedTenantRegisterReq *client.CompletedTenantRegisterRequest, db *sql.DB, c *gin.Context) (*client.UserResponse, error) {
	return nil, nil
}