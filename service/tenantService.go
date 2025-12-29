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
	FindTenantById(ctx context.Context, id int64) (*model.Tenant, error)
	FindAllTenants(ctx context.Context) ([]model.Tenant, error)
	HandleCompleteTenantInfo(ctx context.Context, userId int64, req *client.TenantInfoRequest) (*client.TenantResponse, error)
	FindTenantByUserId(ctx context.Context, userId int64) (*model.Tenant, error)
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
	if err = validateCompleteTenantInfoRequest(ctx, req.Phone, req.Identity, ts.UserRepo); err != nil {
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
