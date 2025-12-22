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

func (tr *TenantRepository) GetTenantByCode(ctx context.Context, code string) (*model.Tenant, error) {
	var err error
	var t model.Tenant
	query := "SELECT * FROM tenants WHERE code = $1 LIMIT 1"
	err = tr.Db.QueryRowContext(ctx, query, code).Scan(
		&t.Id,
		&t.Code,
		&t.Description,
		&t.Type,
		&t.Notes,
		&t.UserId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No tenant found.")
		}
		return nil, err
	}
	return &t, nil
}

func (tr *TenantRepository) GetTenantById(ctx context.Context, id int64) (*model.Tenant, error) {
	var err error
	var t model.Tenant
	query := "SELECT * FROM tenants WHERE id = $1 LIMIT 1"
	err = tr.Db.QueryRowContext(ctx, query, id).Scan(
		&t.Id,
		&t.Code,
		&t.Description,
		&t.Type,
		&t.Notes,
		&t.UserId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No tenant found.")
		}
		return nil, err
	}
	return &t, nil
}

func (tr *TenantRepository) GetTenantByUserId(ctx context.Context, userId int64) (*model.Tenant, error) {
	var err error
	var t model.Tenant
	query := "SELECT * FROM tenants WHERE user_id = $1 LIMIT 1"
	err = tr.Db.QueryRowContext(ctx, query, userId).Scan(
		&t.Id,
		&t.Code,
		&t.Description,
		&t.Type,
		&t.Notes,
		&t.UserId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No tenant found.")
		}
		return nil, err
	}
	return &t, nil
}


func (tr *TenantRepository) ListAllTenants(ctx context.Context) ([]model.Tenant, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM tenants ORDER BY id ASC"
	rows, err = tr.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var tenants []model.Tenant
	for rows.Next() {
		var t model.Tenant
		if err = rows.Scan(
			&t.Id, 
			&t.Code, 
			&t.Description, 
			&t.Type, 
			&t.Notes, 
			&t.UserId, 
		); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	if len(tenants) == 0 {
		return nil, appErr.NotFoundError("No tenant found.")
	}
	return tenants, nil
}
