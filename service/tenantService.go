package service

import (
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
)

type ITenantService interface {
	GetCurrentTenant(c *gin.Context) (*model.Tenant, error)
	GetTenantById(tId int64) (*model.Tenant, error)
	GetListTenant() ([]model.Tenant, error)
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

func (ts *TenantService) GetCurrentTenant(c *gin.Context) (*model.Tenant, error) {
	var err error
	var currUser *model.User

	currUser, err = ts.UserService.GetCurrentUser(c)
	if err != nil {
		return nil, err
	}
	var currTenant *model.Tenant
	currTenant, err = ts.TenantRepo.GetTenantByUserId(currUser.Id)
	if err != nil {
		return nil, err
	}
	return currTenant, nil
}


func (ts *TenantService) GetTenantById(tId int64) (*model.Tenant, error) {
	tenant, err := ts.TenantRepo.GetTenantById(tId)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}


func (ts *TenantService) GetListTenant() ([]model.Tenant, error) {
	tenants, err := ts.TenantRepo.ListAllTenants()
	if err != nil {
		return nil, err
	}
	return tenants, nil
}
