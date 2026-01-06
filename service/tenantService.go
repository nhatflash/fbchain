package service

import (
	"context"
	"time"

	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/security"
	"github.com/shopspring/decimal"
)

type ITenantService interface {
	FindTenantById(ctx context.Context, id int64) (*model.Tenant, error)
	FindAllTenants(ctx context.Context) ([]model.Tenant, error)
	HandleCompleteTenantInfo(ctx context.Context, userId int64, req *client.TenantInfoRequest) (*client.TenantResponse, error)
	FindTenantByUserId(ctx context.Context, userId int64) (*model.Tenant, error)
	HandleCreateNewRestaurantStaff(ctx context.Context, tenantId int64, restaurantId int64, req *client.CreateRestaurantStaffRequest) (*client.RestaurantStaffResponse, error)
}

type TenantService struct {
	TenantRepo 			*repository.TenantRepository
	UserRepo 			*repository.UserRepository
	RestaurantRepo 		*repository.RestaurantRepository
}

func NewTenantService(tr *repository.TenantRepository, ur *repository.UserRepository, rr *repository.RestaurantRepository) ITenantService {
	return &TenantService{
		TenantRepo: tr,
		UserRepo: ur,
	}
}


func (ts *TenantService) FindTenantById(ctx context.Context, id int64) (*model.Tenant, error) {
	tenant, err := ts.TenantRepo.FindTenantById(ctx, id)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}


func (ts *TenantService) FindAllTenants(ctx context.Context) ([]model.Tenant, error) {
	tenants, err := ts.TenantRepo.FindAllTenants(ctx)
	if err != nil {
		return nil, err
	}
	return tenants, nil
}


func (ts *TenantService) HandleCompleteTenantInfo(ctx context.Context, userId int64, req *client.TenantInfoRequest) (*client.TenantResponse, error) {
	var err error
	if err = ts.validateCompleteTenantInfoRequest(ctx, req); err != nil {
		return nil, err
	} 
	
	code := GenerateTenantCode()
	var u *model.User
	var t *model.Tenant
	u, t, err = ts.TenantRepo.CompleteTenantInformation(ctx, req.Phone, req.Identity, req.Address, req.PostalCode, req.ProfileImage, code, req.Description, req.Type, userId)
	if err != nil {
		return nil, err
	}
	return helper.MapToTenantResponse(u, t), nil
}


func (ts *TenantService) FindTenantByUserId(ctx context.Context, userId int64) (*model.Tenant, error) {
	t, err := ts.TenantRepo.FindTenantByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return t, nil
}



func (ts *TenantService) HandleCreateNewRestaurantStaff(ctx context.Context, tenantId int64, restaurantId int64, req *client.CreateRestaurantStaffRequest) (*client.RestaurantStaffResponse, error) {
	var err error
	var r *model.Restaurant
	r, err = ts.RestaurantRepo.FindRestaurantById(ctx, restaurantId)
	if err != nil {
		return nil, err
	}

	if r.TenantId != tenantId {
		return nil, appErr.BadRequestError("The request restaurant does not belong to you.")
	}

	if err = ts.validateCreateNewRestaurantStaffRequest(ctx, req); err != nil {
		return nil, err
	}

	var passwordHash string
	passwordHash, err = security.GenerateHashedPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var birthdate *time.Time
	birthdate, err = helper.ConvertToDate(req.Birthdate)
	if err != nil {
		return nil, err
	}
	
	var salary decimal.Decimal
	salary, err = decimal.NewFromString(req.Salary)
	if err != nil {
		return nil, appErr.BadRequestError("Invalid salary format.")
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
	
	code := GenerateRestaurantStaffCode()

	var u *model.User
	var rs *model.RestaurantStaff
	u, rs, err = ts.UserRepo.CreateRestaurantStaffUser(ctx, req.Email, passwordHash, req.Phone, req.Identity, req.FirstName, req.LastName, &req.Gender, birthdate, req.PostalCode, req.Address, req.ProfileImage, restaurantId, code, &req.Type, shiftStart, shiftEnd, salary, req.Notes)
	if err != nil {
		return nil, err
	}

	return helper.MapToRestaurantStaffResponse(u, rs), nil
}


func (ts *TenantService) validateCompleteTenantInfoRequest(ctx context.Context, req *client.TenantInfoRequest) error {
	var err error
	var exist bool
	exist, err = ts.UserRepo.CheckUserPhoneExists(ctx, req.Phone)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("This phone number is already registered.")
	}

	exist, err = ts.UserRepo.CheckUserIdentityExists(ctx, req.Identity)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("This identity is already registered.")
	}
	return nil
}


func (ts *TenantService) validateCreateNewRestaurantStaffRequest(ctx context.Context, req *client.CreateRestaurantStaffRequest) error {
	var err error
	var exist bool
	exist, err = ts.UserRepo.CheckUserEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("This email is already in use.")
	}

	exist, err = ts.UserRepo.CheckUserPhoneExists(ctx, req.Phone)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("This phone number is already in use.")
	}

	exist, err = ts.UserRepo.CheckUserIdentityExists(ctx, req.Identity)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("This identity number is already in use.")
	}

	if req.ConfirmPassword != req.Password {
		return appErr.BadRequestError("Confirm password does not match.")
	}
	return nil
}


