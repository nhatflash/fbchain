package repository

import (
	"database/sql"

	"github.com/nhatflash/fbchain/enum"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
)

func CreateTenantInformation(code string, description string, tenantType *enum.TenantType, userId int64, db *sql.DB) (*model.Tenant, error) {
	var err error
	_, err = db.Exec("INSERT INTO tenants (code, description, type, user_id) VALUES ($1, $2, $3, $4)", code, description, tenantType, userId)
	if err != nil {
		return nil, err
	}

	var tenant *model.Tenant
	tenant, err = GetTenantByCode(code, db)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func GetTenantByCode(code string, db *sql.DB) (*model.Tenant, error) {
	var rows *sql.Rows
	var err error
	rows, err = db.Query("SELECT * FROM tenants WHERE code = $1 LIMIT 1", code)
	if err != nil {
		return nil, err
	}
	var tenants []model.Tenant
	for rows.Next() {
		var t model.Tenant
		err = rows.Scan(&t.Id, &t.Code, &t.Description, &t.Type, &t.Notes, &t.UserId)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	if len(tenants) == 0 {
		return nil, appErr.NotFoundError("No tenant found")
	}
	return &tenants[0], nil
}

func GetTenantById(tId int64, db *sql.DB) (*model.Tenant, error) {
	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT * FROM tenants WHERE id = $1 LIMIT 1", tId)
	if err != nil {
		return nil, err
	}
	var tenants []model.Tenant
	for rows.Next() {
		var tenant model.Tenant
		err = rows.Scan(&tenant.Id, &tenant.Code, &tenant.Description, &tenant.Type, &tenant.Notes, &tenant.UserId)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}
	if len(tenants) == 0 {
		return nil, appErr.NotFoundError("No tenant found.")
	}
	return &tenants[0], nil
}

func GetTenantByUserId(uId int64, db *sql.DB) (*model.Tenant, error) {
	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT * FROM tenants WHERE user_id = $1 LIMIT 1", uId)
	if err != nil {
		return nil, err
	}
	var tenants []model.Tenant
	for rows.Next() {
		var tenant model.Tenant
		err = rows.Scan(&tenant.Id, &tenant.Code, &tenant.Description, &tenant.Type, &tenant.Notes, &tenant.UserId)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}
	if len(tenants) == 0 {
		return nil, appErr.NotFoundError("No tenant found.")
	}
	return &tenants[0], nil
}
