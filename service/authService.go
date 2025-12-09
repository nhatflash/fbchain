package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/helper"
	appError "github.com/nhatflash/fbchain/error"
)

func HandleLogin(loginReq *client.LoginRequest, db *sql.DB, c *gin.Context) (string, error) {
	user, userErr := repository.GetSignInUser(loginReq.Login, loginReq.Password, db)
	if userErr != nil {
		return "", userErr
	}
	return user.Email, nil
}

func HandleRegisterTenant(registerTenantReq *client.RegisterTenantRequest, db *sql.DB, c *gin.Context) (*client.UserResponse, error) {
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