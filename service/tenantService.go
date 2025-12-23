package service

import (
	"context"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/helper"
)

type ITenantService interface {
	GetCurrentTenant(ctx context.Context) (*model.Tenant, error)
	GetTenantById(ctx context.Context, id int64) (*model.Tenant, error)
	GetListTenant(ctx context.Context) ([]model.Tenant, error)
	HandleCompleteTenantInfo(ctx context.Context, req *client.TenantInfoRequest) (*client.TenantResponse, error)
}

type TenantService struct {
	TenantRepo 			*repository.TenantRepository
	UserService			IUserService
	UserRepo 			*repository.UserRepository
}

func NewTenantService(tr *repository.TenantRepository, us IUserService, ur *repository.UserRepository) ITenantService {
	return &TenantService{
		TenantRepo: tr,
		UserService: us,
		UserRepo: ur,
	}
}

func (ts *TenantService) GetCurrentTenant(ctx context.Context) (*model.Tenant, error) {
	var err error
	var currUser *model.User

	currUser, err = ts.UserService.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	var currTenant *model.Tenant
	currTenant, err = ts.TenantRepo.GetTenantByUserId(ctx, currUser.Id)
	if err != nil {
		return nil, err
	}
	return currTenant, nil
}


func (ts *TenantService) GetTenantById(ctx context.Context, id int64) (*model.Tenant, error) {
	tenant, err := ts.TenantRepo.GetTenantById(ctx, id)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}


func (ts *TenantService) GetListTenant(ctx context.Context) ([]model.Tenant, error) {
	tenants, err := ts.TenantRepo.ListAllTenants(ctx)
	if err != nil {
		return nil, err
	}
	return tenants, nil
}


func (ts *TenantService) HandleCompleteTenantInfo(ctx context.Context, req *client.TenantInfoRequest) (*client.TenantResponse, error) {
	phone := req.Phone
	identity := req.Identity
	address := req.Address
	postalCode := req.PostalCode
	description := req.Description
	tenantType := req.Type
	profileImage := req.ProfileImage

	var err error
	var currUser *model.User
	currUser, err = ts.UserService.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if err = ValidateCompleteTenantInfoRequest(ctx, phone, identity, currUser, ts.UserRepo); err != nil {
		return nil, err
	} 
	
	code := GenerateTenantCode()
	var tenant *model.Tenant
	userId := currUser.Id
	currUser, tenant, err = ts.TenantRepo.CompleteTenantInformation(ctx, phone, identity, address, postalCode, profileImage, code, description, tenantType, userId)
	if err != nil {
		return nil, err
	}
	return helper.MapToTenantResponse(currUser, tenant), nil
}


func ValidateCompleteTenantInfoRequest(ctx context.Context, phone string, identity string, u *model.User, ur *repository.UserRepository) error {
	if u.IsVerified {
		return appErr.BadRequestError("This tenant is already verified.")
	}
	var err error
	var exist bool
	exist, err = ur.CheckUserPhoneExists(ctx, phone)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("This phone number is already registered.")
	}

	exist, err = ur.CheckUserIdentityExists(ctx, identity)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("This identity is already registered.")
	}
	return nil
}
