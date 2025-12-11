package repository

import (
	"database/sql"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/enum"
	"github.com/nhatflash/fbchain/model"
)

func CreateTenantInformation(code string, description string, tenantType *enum.TenantType, userId int64, db *sql.DB) (*model.Tenant, error) {

	_, insertErr := db.Exec("INSERT INTO tenants (code, description, type, user_id) VALUES ($1, $2, $3, $4)", code, description, tenantType, userId)
	if insertErr != nil {
		return nil, insertErr
	}
	tenant, err := GetTenantByCode(code, db)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func GetTenantByCode(code string, db *sql.DB) (*model.Tenant, error) {
	rows, dbErr := db.Query("SELECT * FROM tenants WHERE code = $1 LIMIT 1", code)
	if dbErr != nil {
		return nil, dbErr
	}
	var tenants []model.Tenant
	for rows.Next() {
		var tenant model.Tenant
		scanErr := rows.Scan(&tenant.Id, &tenant.Code, &tenant.Description, &tenant.Type, &tenant.Notes, &tenant.UserId)
		if scanErr != nil {
			return nil, scanErr
		}
		tenants = append(tenants, tenant)
	}
	if len(tenants) == 0 {
		return nil, appErr.NotFoundError("No tenant found")
	}
	return &tenants[0], nil
}