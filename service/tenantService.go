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
	GetTenantById(ctx context.Context, id int64) (*model.Tenant, error)
	GetListTenant(ctx context.Context) ([]model.Tenant, error)
	HandleCompleteTenantInfo(ctx context.Context, userId int64, req *client.TenantInfoRequest) (*client.TenantResponse, error)
	GetTenantByUserId(ctx context.Context, userId int64) (*model.Tenant, error)
}

type TenantService struct {
	TenantRepo 			*repository.TenantRepository
	UserRepo 			*repository.UserRepository
}

func NewTenantService(tr *repository.TenantRepository, ur *repository.UserRepository) ITenantService {
	return &TenantService{
		TenantRepo: tr,
		UserRepo: ur,
	}
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


func (ts *TenantService) HandleCompleteTenantInfo(ctx context.Context, userId int64, req *client.TenantInfoRequest) (*client.TenantResponse, error) {
	phone := req.Phone
	identity := req.Identity
	address := req.Address
	postalCode := req.PostalCode
	description := req.Description
	tenantType := req.Type
	profileImage := req.ProfileImage

	var err error
	if err = validateCompleteTenantInfoRequest(ctx, phone, identity, ts.UserRepo); err != nil {
		return nil, err
	} 
	
	code := GenerateTenantCode()
	var u *model.User
	var t *model.Tenant
	u, t, err = ts.TenantRepo.CompleteTenantInformation(ctx, phone, identity, address, postalCode, profileImage, code, description, tenantType, userId)
	if err != nil {
		return nil, err
	}
	return helper.MapToTenantResponse(u, t), nil
}


func (ts *TenantService) GetTenantByUserId(ctx context.Context, userId int64) (*model.Tenant, error) {
	t, err := ts.TenantRepo.GetTenantByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return t, nil
}


func validateCompleteTenantInfoRequest(ctx context.Context, phone string, identity string, ur *repository.UserRepository) error {
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
