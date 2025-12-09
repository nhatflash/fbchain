package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/helper"
)

func HandleLogin(loginReq *client.LoginRequest, db *sql.DB, c *gin.Context) (string, error) {
	user := repository.GetSignInUser(loginReq.Login, loginReq.Password, db)
	if user == nil {
		return "", errors.New("invalid credential")
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
		return nil, errors.New("user with this email already exist: " + email)
	}
	if (password != confirmPassword) {
		return nil, errors.New("confirm password does not match")
	}
	birthdate, dateErr := helper.ConvertToDate(birthdateStr)
	if dateErr != nil {
		return nil, errors.New(dateErr.Error())
	}
	newTenant, dbErr := repository.RegisterTenant(email, firstName, lastName, password, &gender, birthdate, db)
	if dbErr != nil {
		return nil, errors.New(dbErr.Error())
	}
	return helper.MapToUserResponse(newTenant), nil
}