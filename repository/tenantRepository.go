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

func (tr *TenantRepository) CompleteTenantInformation(ctx context.Context, phone string, identity string, address string, postalCode string, profileImage *string, code string, description *string, tenantType *enum.TenantType, userId int64) (*model.User, *model.Tenant, error) {
	var err error
	var tx *sql.Tx
	tx, err = tr.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	defer tx.Rollback()

	var u model.User
	userQuery := "UPDATE users SET phone = $1, identity = $2, address = $3, postal_code = $4, profile_image = $5, is_verified = $6 WHERE id = $7 RETURNING *"
	if err = tx.QueryRowContext(ctx, userQuery, phone, identity, address, postalCode, profileImage, true, userId).Scan(
		&u.Id,
		&u.Email,
		&u.Password,
		&u.Role, 
		&u.Phone,
		&u.Identity,
		&u.FirstName,
		&u.LastName,
		&u.Gender,
		&u.Birthdate,
		&u.PostalCode, 
		&u.Address,
		&u.ProfileImage,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsVerified,
	); err != nil {
		return nil, nil, err
	}

	var t model.Tenant
	tenantQuery := "INSERT INTO tenants (code, description, type, user_id) VALUES ($1, $2, $3, $4) RETURNING *"
	if err = tx.QueryRowContext(ctx, tenantQuery, code, description, tenantType, userId).Scan(
		&t.Id,
		&t.Code,
		&t.Description,
		&t.Type,
		&t.Notes,
		&t.UserId,
	); err != nil {
		return nil, nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, err
	}

	return &u, &t, nil
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
