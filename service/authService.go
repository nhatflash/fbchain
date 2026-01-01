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
	var err error
	var foundUser *model.User
	if strings.Contains(req.Login, "@") {
		foundUser, err = as.UserRepo.FindUserByEmail(ctx, req.Login)
		if err != nil {
			return nil, appErr.UnauthorizedError("Invalid credentials.")
		}
	} else {
		foundUser, err = as.UserRepo.FindUserByPhone(ctx, req.Login)
		if err != nil {
			return nil, appErr.UnauthorizedError("Invalid credentials.")
		}
	}
	if !security.VerifyPassword(req.Password, foundUser.Password) {
		return nil, appErr.UnauthorizedError("Invalid credentials.")
	}
	
	var accessToken string
	accessToken, err = security.GenerateJwtAccessToken(foundUser)
	if err != nil {
		return nil, err
	}
	var refreshToken string
	refreshToken, err = security.GenerateJwtRefreshToken(foundUser)
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
	var err error
	if err = ValidateSignUpRequest(ctx, req.Email, req.Password, req.ConfirmPassword, as.UserRepo); err != nil {
		return nil, err
	}

	var birthdate *time.Time
	birthdate, err = helper.ConvertToDate(req.Birthdate)
	if err != nil {
		return nil, err
	}
	var hashedPassword string
	hashedPassword, err = security.GenerateHashedPassword(req.Password)
	if err != nil {
		return nil, err
	}
	var tenantUser *model.User
	tenantUser, err = as.UserRepo.CreateTenantUser(ctx, req.FirstName, req.LastName, req.Email, hashedPassword, birthdate, req.Gender)
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
	validTimeKey := constant.UserChangePasswordTimeKey + userId
	var exists int64
	exists, err  = as.Rdb.Exists(ctx, validTimeKey).Result()
	if err != nil {
		return "", err
	}
	if exists == 1 {
		return "", appErr.UnauthorizedError("You are on a password change session. Please try again later.")
	}

	otpLen := 6
	var otp string
	otp, err = security.GenerateOTPCode(otpLen)
	if err != nil {
		return "", err
	}
	duration := time.Duration(15) * time.Minute
	otpKey := constant.UserVerifyPasswordOTPKey + userId
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
	otpKey := constant.UserVerifyPasswordOTPKey + userId
	actualOTP, err = as.Rdb.Get(ctx, otpKey).Result()
	if err == redis.Nil {
		return appErr.UnauthorizedError("OTP code expired or not found.")
	}
	if err != nil {
		return err
	}

	if security.VerifyOTPCode(req.VerifiedCode, actualOTP) {
		as.Rdb.Del(ctx, otpKey)
		validTimeKey := constant.UserChangePasswordTimeKey + userId
		duration := time.Duration(15) * time.Minute
		if err = as.Rdb.Set(ctx, validTimeKey, "true", duration).Err(); err != nil {
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
	validTimeKey := constant.UserChangePasswordTimeKey + userId

	var exists int64
	exists, err = as.Rdb.Exists(ctx, validTimeKey).Result()
	if err != nil {
		return err
	}
	if exists != 1 {
		return appErr.UnauthorizedError("Password change session has expired.")
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



