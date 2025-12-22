package repository

import (
	"context"
	"database/sql"

	"github.com/nhatflash/fbchain/enum"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
)

type TenantRepository struct {
	Db 		*sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{
		Db: db,
	}
}

func (tr *TenantRepository) CreateTenantInformation(ctx context.Context, code string, description string, tenantType *enum.TenantType, userId int64) (*model.Tenant, error) {
	var err error
	var tx *sql.Tx
	tx, err = tr.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var t model.Tenant
	query := "INSERT INTO tenants (code, description, type, user_id) VALUES ($1, $2, $3, $4) RETURNING *"
	if err = tx.QueryRowContext(ctx, query, code, description, tenantType, userId).Scan(
		&t.Id,
		&t.Code,
		&t.Description,
		&t.Type,
		&t.Notes,
		&t.UserId,
	); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &t, nil
}

func (tr *TenantRepository) GetTenantByCode(code string) (*model.Tenant, error) {
	var rows *sql.Rows
	var err error
	rows, err = tr.Db.Query("SELECT * FROM tenants WHERE code = $1 LIMIT 1", code)
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

func (tr *TenantRepository) GetTenantById(tId int64) (*model.Tenant, error) {
	var rows *sql.Rows
	var err error

	rows, err = tr.Db.Query("SELECT * FROM tenants WHERE id = $1 LIMIT 1", tId)
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

func (tr *TenantRepository) GetTenantByUserId(uId int64) (*model.Tenant, error) {
	var rows *sql.Rows
	var err error

	rows, err = tr.Db.Query("SELECT * FROM tenants WHERE user_id = $1 LIMIT 1", uId)
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


func (tr *TenantRepository) ListAllTenants() ([]model.Tenant, error) {
	var rows *sql.Rows
	var err error

	rows, err = tr.Db.Query("SELECT * FROM tenants ORDER BY id ASC")
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
	return tenants, nil
}
