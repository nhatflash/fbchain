package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
)

func GenerateTenantCode() string {
	now := time.Now()
	unixMilli := now.UnixMilli()
	return fmt.Sprintf("TENANT-%d", unixMilli)
}

func GetCurrentTenant(c *gin.Context, db *sql.DB) (*model.Tenant, error) {
	var err error
	var currUser *model.User

	currUser, err = GetCurrentUser(c, db)
	if err != nil {
		return nil, err
	}
	var currTenant *model.Tenant
	currTenant, err = repository.GetTenantByUserId(currUser.Id, db)
	if err != nil {
		return nil, err
	}
	return currTenant, nil
}
