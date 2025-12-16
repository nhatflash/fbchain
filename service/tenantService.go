package service

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
)

type ITenantService interface {
	GetCurrentTenant(c *gin.Context) (*model.Tenant, error)
	GetTenantById(tId int64) (*model.Tenant, error)
}

type TenantService struct {
	Db *sql.DB
}

// GetTenantById implements [ITenantService].
func (t *TenantService) GetTenantById(tId int64) (*model.Tenant, error) {
	panic("unimplemented")
}

func NewTenantService(db *sql.DB) ITenantService {
	return &TenantService{
		Db: db,
	}
}

func (t *TenantService) GetCurrentTenant(c *gin.Context) (*model.Tenant, error) {
	var err error
	var currUser *model.User

	userService := NewUserService(t.Db)
	currUser, err = userService.GetCurrentUser(c)
	if err != nil {
		return nil, err
	}
	var currTenant *model.Tenant
	currTenant, err = repository.GetTenantByUserId(currUser.Id, t.Db)
	if err != nil {
		return nil, err
	}
	return currTenant, nil
}
