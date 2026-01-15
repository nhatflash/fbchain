package service

import (
	"context"
	"time"

	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/enum"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/security"
	"github.com/shopspring/decimal"
)

type AdminService struct {
	UserRepo 			*repository.UserRepository
}


type IAdminService interface {
	HandleCreateNewStaff(ctx context.Context, req *client.CreateStaffRequest) (*client.StaffResponse, error)
}


func NewAdminService(ur *repository.UserRepository) IAdminService {
	return &AdminService{
		UserRepo: ur,
	}
}


func (as *AdminService) HandleCreateNewStaff(ctx context.Context, req *client.CreateStaffRequest) (*client.StaffResponse, error) {
	var err error
	if err = as.validateCreateStaffRequest(ctx, req); err != nil {
		return nil, err
	}
	var birthdate *time.Time
	birthdate, err = helper.ConvertToDate(req.Birthdate)
	if err != nil {
		return nil, err
	}

	var shiftStart *time.Time
	if req.ShiftStart != nil {
		shiftStart, err = helper.ConvertToTime(*req.ShiftStart)
		if err != nil {
			return nil, err
		}
	}
	var shiftEnd *time.Time
	if req.ShiftEnd != nil {
		shiftEnd, err = helper.ConvertToTime(*req.ShiftEnd)
		if err != nil {
			return nil, err
		}
	}

	var salary decimal.Decimal
	salary, err = decimal.NewFromString(req.Salary)
	if err != nil {
		return nil, appErr.BadRequestError("Invalid salary format.")
	}
	
	var role enum.Role
	var code string
	if req.Role == client.SR_MANAGER {
		role = enum.ROLE_MANAGER
		code = GenerateManagerCode()
	} else {
		role = enum.ROLE_STAFF
		code = GenerateStaffCode()
	}

	var passwordHash string
	passwordHash, err = security.GenerateHashedPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var u *model.User
	var s *model.Staff
	u, s, err = as.UserRepo.CreateStaffUser(ctx, req.Email, passwordHash, &role, req.Phone, req.Identity, req.FirstName, req.LastName, &req.Gender, birthdate, req.PostalCode, req.Address, req.ProfileImage, code, &req.ShiftType, shiftStart, shiftEnd, salary, req.Notes)
	if err != nil {
		return nil, err
	}
	return helper.MapToStaffResponse(u, s), nil
} 

func (as *AdminService) validateCreateStaffRequest(ctx context.Context, req *client.CreateStaffRequest) error {
	var err error
	var exist bool
	exist, err = as.UserRepo.CheckUserEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("This email is already in use.")
	}

	exist, err = as.UserRepo.CheckUserPhoneExists(ctx, req.Phone)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("This phone number is already in use.")
	}

	exist, err = as.UserRepo.CheckUserIdentityExists(ctx, req.Identity)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("This identity number is already in use.")
	}

	if req.Password != req.ConfirmPassword {
		return appErr.BadRequestError("Confirm password does not match.")
	}

	if req.ShiftType == enum.STAFF_PARTTIME && (req.ShiftStart == nil || req.ShiftEnd == nil) {
		return appErr.BadRequestError("Staff shift start and shift end must be specified if the staff is working part time.")
	}

	return nil

}
