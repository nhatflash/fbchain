package service

import (
	"context"

	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
)

type ITenantService interface {
	GetCurrentTenant(ctx context.Context) (*model.Tenant, error)
	GetTenantById(ctx context.Context, id int64) (*model.Tenant, error)
	GetListTenant(ctx context.Context) ([]model.Tenant, error)
}

type TenantService struct {
	TenantRepo 			*repository.TenantRepository
	UserService			IUserService
}

func NewTenantService(tr *repository.TenantRepository, us IUserService) ITenantService {
	return &TenantService{
		TenantRepo: tr,
		UserService: us,
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
