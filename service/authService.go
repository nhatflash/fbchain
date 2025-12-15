package service

import (
	"database/sql"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/helper"
	appError "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/security"
	"time"
	"strings"
	"github.com/nhatflash/fbchain/model"
)

type AuthService interface {
	HandleSignIn(signInReq *client.SignInRequest, db *sql.DB) (*client.SignInResponse, error)
	HandleTenantSignUp(tenantSignUpReq *client.TenantSignUpRequest, db *sql.DB) (*client.TenantResponse, error)
}


// Sign in method
func HandleSignIn(signInReq *client.SignInRequest, db *sql.DB) (*client.SignInResponse, error) {
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

// tenant sign up method
func HandleTenantSignUp(tenantSignUpReq *client.TenantSignUpRequest, db *sql.DB) (*client.TenantResponse, error) {
	firstName := tenantSignUpReq.FirstName
	lastName := tenantSignUpReq.LastName
	email := tenantSignUpReq.Email
	password := tenantSignUpReq.Password
	confirmPassword := tenantSignUpReq.Password
	birthdateStr := tenantSignUpReq.Birthdate
	gender := tenantSignUpReq.Gender
	phone := tenantSignUpReq.Phone
	identity := tenantSignUpReq.Identity
	address := tenantSignUpReq.Address
	postalCode := tenantSignUpReq.PostalCode
	profileImage := tenantSignUpReq.ProfileImage
	description := tenantSignUpReq.Description
	tenantType := tenantSignUpReq.Type

	validateErr := validateSignUpRequest(email, phone, identity, password, confirmPassword, db)

	if validateErr != nil {
		return nil, validateErr
	}

	birthdate, dateErr := helper.ConvertToDate(birthdateStr)
	if dateErr != nil {
		return nil, dateErr
	}
	hashedPassword, hashErr := security.GenerateHashedPassword(password)
	if hashErr != nil {
		return nil, appError.InternalError("An unexpected error when creating hash password")
	}
	
	userTenant, userErr := repository.CreateTenantUser(firstName, lastName, email, hashedPassword, birthdate, gender, phone, identity, address, postalCode, profileImage, db)
	if userErr != nil {
		return nil, userErr
	}
	code := GenerateTenantCode()
	tenant, tenantErr := repository.CreateTenantInformation(code, description, tenantType, userTenant.Id, db)

	if tenantErr != nil {
		return nil, tenantErr
	}
	return helper.MapToTenantResponse(userTenant, tenant), nil
}


func validateSignUpRequest(email string, phone string, identity string, password string, confirmPassword string, db *sql.DB) error {
	if repository.CheckUserEmailExists(email, db) {
		return appError.BadRequestError("User with this email already exists.")
	}
	if repository.CheckUserPhoneExists(phone, db) {
		return appError.BadRequestError("User with this phone already exists")
	}
	if repository.CheckUserIdentityExists(identity, db) {
		return appError.BadRequestError("User with this identity already exists")
	}
	if password != confirmPassword {
		return appError.BadRequestError("Confirm password does not match.")
	}
	return nil
}