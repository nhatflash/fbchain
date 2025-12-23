package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nhatflash/fbchain/constant"

	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/middleware"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/security"
	"github.com/redis/go-redis/v9"
)

type IAuthService interface {
	HandleSignIn(ctx context.Context, req *client.SignInRequest) (*client.SignInResponse, error)
	HandleTenantUserSignUp(ctx context.Context, req *client.TenantSignUpRequest) (*client.UserResponse, error)
	GenerateChangePasswordVerifyOTP(ctx context.Context) (string, error)
	HandleVerifyChangePassword(ctx context.Context, req *client.VerifyChangePasswordRequest) error
	HandleChangePassword(ctx context.Context, req *client.ChangePasswordRequest) (error)
}

type AuthService struct {
	UserRepo 		*repository.UserRepository
	TenantRepo 		*repository.TenantRepository
	Rdb 			*redis.Client
}

func NewAuthService(ur *repository.UserRepository, tr *repository.TenantRepository, rdb *redis.Client) IAuthService {
	return &AuthService{
		UserRepo: ur,
		TenantRepo: tr,
		Rdb: rdb,
	}
}


// Sign in method
func (as *AuthService) HandleSignIn(ctx context.Context, req *client.SignInRequest) (*client.SignInResponse, error) {
	login := req.Login

	var loggedUser *model.User
	var err error
	var foundUser *model.User
	if strings.Contains(login, "@") {
		foundUser, err = as.UserRepo.GetUserByEmail(ctx, req.Login)
		if err != nil {
			return nil, appErr.UnauthorizedError("Invalid credentials.")
		}
	} else {
		foundUser, err = as.UserRepo.GetUserByPhone(ctx, req.Login)
		if err != nil {
			return nil, appErr.UnauthorizedError("Invalid credentials.")
		}
	}
	loggedUser = foundUser
	if !security.VerifyPassword(req.Password, loggedUser.Password) {
		return nil, appErr.UnauthorizedError("Invalid credentials.")
	}
	
	var accessToken string
	accessToken, err = security.GenerateJwtAccessToken(loggedUser)
	if err != nil {
		return nil, err
	}
	var refreshToken string
	refreshToken, err = security.GenerateJwtRefreshToken(loggedUser)
	if err != nil {
		return nil, err
	}
	res := &client.SignInResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		LastLogin: time.Now(),
	}
	return res, nil
}

// tenant sign up method
func (as *AuthService) HandleTenantUserSignUp(ctx context.Context, req *client.TenantSignUpRequest) (*client.UserResponse, error) {
	firstName := req.FirstName
	lastName := req.LastName
	email := req.Email
	password := req.Password
	confirmPassword := req.Password
	birthdateStr := req.Birthdate
	gender := req.Gender

	var err error
	if err = ValidateSignUpRequest(ctx, email, password, confirmPassword, as.UserRepo); err != nil {
		return nil, err
	}

	var birthdate *time.Time
	birthdate, err = helper.ConvertToDate(birthdateStr)
	if err != nil {
		return nil, err
	}
	var hashedPassword string
	hashedPassword, err = security.GenerateHashedPassword(password)
	if err != nil {
		return nil, err
	}
	var tenantUser *model.User
	tenantUser, err = as.UserRepo.CreateTenantUser(ctx, firstName, lastName, email, hashedPassword, birthdate, gender)
	if err != nil {
		return nil, err
	}
	return helper.MapToUserResponse(tenantUser), nil
}


func (as *AuthService) GenerateChangePasswordVerifyOTP(ctx context.Context) (string, error) {
	var err error
	var claims *security.JwtAccessClaims
	claims, err = GetCurrentClaims(ctx)
	if err != nil {
		return "", err
	}
	userId := strconv.FormatInt(claims.UserId, 10)
	validTimeKey := constant.USER_CHANGE_PASSWORD_TIME_KEY + userId
	var validTimeValue string
	validTimeValue, _ = as.Rdb.Get(ctx, validTimeKey).Result()
	if validTimeValue == constant.USER_CHANGE_PASSWORD_TIME_VALUE {
		return "", appErr.BadRequestError("You are already on an attempt for password changing. Please try again later.")
	}

	otpLen := 6
	var otp string
	otp, err = security.GenerateOTPCode(otpLen)
	if err != nil {
		return "", err
	}
	duration := time.Duration(constant.VERIFY_PASSWORD_OTP_EXPIRATION_MIN) * time.Minute
	otpKey := constant.USER_VERIFY_PASSWORD_OTP_KEY + userId
	err = as.Rdb.Set(ctx, otpKey, otp, duration).Err()
	if err != nil {
		return "", err
	}
	return otp, nil
}


func (as *AuthService) HandleVerifyChangePassword(ctx context.Context, req *client.VerifyChangePasswordRequest) error {
	var err error
	var claims *security.JwtAccessClaims
	claims, err = GetCurrentClaims(ctx)
	if err != nil {
		return err
	}
	userId := strconv.FormatInt(claims.UserId, 10)
	var actualOTP string
	otpKey := constant.USER_VERIFY_PASSWORD_OTP_KEY + userId
	actualOTP, err = as.Rdb.Get(ctx, otpKey).Result()
	if err == redis.Nil {
		return appErr.UnauthorizedError("OTP code expired or not found.")
	}
	if err != nil {
		return err
	}

	if security.VerifyOTPCode(req.VerifiedCode, actualOTP) {
		as.Rdb.Del(ctx, otpKey)
		validTimeKey := constant.USER_CHANGE_PASSWORD_TIME_KEY + userId
		duration := time.Duration(constant.CHANGE_PASSWORD_TIME) * time.Minute
		if err = as.Rdb.Set(ctx, validTimeKey, constant.USER_CHANGE_PASSWORD_TIME_VALUE, duration).Err(); err != nil {
			return err
		}
		return nil
	}
	return appErr.UnauthorizedError("Your otp is not valid or expired.")
}


func (as *AuthService) HandleChangePassword(ctx context.Context, req *client.ChangePasswordRequest) (error) {
	var err error
	var claims *security.JwtAccessClaims
	claims, err = GetCurrentClaims(ctx)
	if err != nil {
		return err
	}
	userId := strconv.FormatInt(claims.UserId, 10)
	validTimeKey := constant.USER_CHANGE_PASSWORD_TIME_KEY + userId
	_, err = as.Rdb.Get(ctx, validTimeKey).Result()
	if err == redis.Nil {
		return appErr.UnauthorizedError("Password change session expired. Please perform another request.")
	}
	if err != nil {
		return err
	}

	newPassword := req.NewPassword
	confirmedNewPassword := req.ConfirmNewPassword

	if !IsConfirmedPasswordMatches(newPassword, confirmedNewPassword) {
		return appErr.BadRequestError("Confirmed password does not match.")
	}
	var hashedPassword string
	hashedPassword, err = security.GenerateHashedPassword(newPassword)
	if err != nil {
		return err
	}

	if err = as.UserRepo.ChangeUserPassword(ctx, claims.UserId, hashedPassword); err != nil {
		return err
	}
	as.Rdb.Del(ctx, validTimeKey)
	return nil
}



func GetCurrentClaims(ctx context.Context) (*security.JwtAccessClaims, error) {
	claims, ok := ctx.Value(middleware.UserKey{}).(*security.JwtAccessClaims)
	if !ok || claims == nil {
		return nil, appErr.UnauthorizedError("Authentication is required.")
	}
	return claims, nil
}



func GenerateTenantCode() string {
	now := time.Now()
	unixMilli := now.UnixMilli()
	return fmt.Sprintf("TENANT-%d", unixMilli)
}


func GenerateStaffCode() string {
	now := time.Now()
	unixMilli := now.UnixMilli()
	return fmt.Sprintf("STAFF-%d", unixMilli)
}

func ValidateSignUpRequest(ctx context.Context, email string, password string, confirmPassword string, ur *repository.UserRepository) error {
	var err error
	var exist bool
	exist, err = ur.CheckUserEmailExists(ctx, email)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("User with this email already exists.")
	}
	if !IsConfirmedPasswordMatches(password, confirmPassword) {
		return appErr.BadRequestError("Confirm password does not match.")
	}
	return nil
}


func IsConfirmedPasswordMatches(password string, confirmedPassword string) bool {
	return password == confirmedPassword
}



