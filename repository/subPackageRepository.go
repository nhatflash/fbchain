package repository

import (
	"context"
	"database/sql"

	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
	"github.com/shopspring/decimal"
)

type SubPackageRepository struct {
	Db			*sql.DB
}

func NewSubPackageRepository(db *sql.DB) *SubPackageRepository {
	return &SubPackageRepository{
		Db: db,
	}
}

func (spr *SubPackageRepository) CheckSubPackageNameExists(ctx context.Context, name string) (bool, error) {
	query := "SELECT name FROM sub_packages WHERE name = $1 LIMIT 1"
	rows, err := spr.Db.QueryContext(ctx, query, name)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (spr *SubPackageRepository) CreateNewSubPackage(ctx context.Context, name string, description *string, durationMonth int, price decimal.Decimal, image *string) (*model.SubPackage, error) {
	var err error
	var tx *sql.Tx
	tx, err = spr.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := "INSERT INTO sub_packages (name, description, duration_month, price, is_active, image) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *"
	var s model.SubPackage
	if err = tx.QueryRowContext(ctx, query, name, description, durationMonth, price, true, image).Scan(
		&s.Id,
		&s.Name,
		&s.Description,
		&s.DurationMonth,
		&s.Price,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.Image,
	); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &s, nil
}

func (spr *SubPackageRepository) FindSubPackageByName(ctx context.Context, name string) (*model.SubPackage, error) {
	var s model.SubPackage
	query := "SELECT * FROM sub_packages WHERE name = $1 LIMIT 1" 
	err := spr.Db.QueryRowContext(ctx, query, name).Scan(
		&s.Id,
		&s.Name,
		&s.Description,
		&s.DurationMonth,
		&s.Price,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.Image,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No subscription package found.")
		}
		return nil, err
	}
	return &s, nil
}

func (spr *SubPackageRepository) AnySubPackageExists(ctx context.Context) (bool, error) {
	rows, err := spr.Db.QueryContext(ctx, "SELECT id FROM sub_packages")
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (spr *SubPackageRepository) IsSubPackageExist(ctx context.Context, sId int64) (bool, error) {
	query := "SELECT id FROM sub_packages WHERE id = $1 LIMIT 1"
	rows, err := spr.Db.QueryContext(ctx, query, sId)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (spr *SubPackageRepository) FindSubPackageById(ctx context.Context, id int64) (*model.SubPackage, error) {
	var s model.SubPackage
	query := "SELECT * FROM sub_packages WHERE id = $1 LIMIT 1"
	err := spr.Db.QueryRowContext(ctx, query, id).Scan(
		&s.Id,
		&s.Name,
		&s.Description, 
		&s.DurationMonth,
		&s.Price, 
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.Image,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No subscription package found.")
		}
		return nil, err
	}
	return &s, nil
}

func (spr *SubPackageRepository) FindFirstSubPackage(ctx context.Context) (*model.SubPackage, error) {
	var s model.SubPackage
	query := "SELECT * FROM sub_packages ORDER BY id ASC LIMIT 1"
	err := spr.Db.QueryRowContext(ctx, query).Scan(
		&s.Id,
		&s.Name,
		&s.Description,
		&s.DurationMonth,
		&s.Price,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.Image,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No subscription package found.")
		}
	}
	return &s, nil
}


func (spr *SubPackageRepository) FindAllSubPackages(ctx context.Context) ([]model.SubPackage, error) {
	var err error
	var rows *sql.Rows
	query := "SELECT * FROM sub_packages"
	rows, err = spr.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var subPackages []model.SubPackage
	for rows.Next() {
		var s model.SubPackage
		if err = rows.Scan(
			&s.Id,
			&s.Name,
			&s.Description,
			&s.DurationMonth,
			&s.Price,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.Image,
		); err != nil {
			return nil, err
		}
		subPackages = append(subPackages, s)
	}
	if len(subPackages) == 0 {
		return nil, appErr.NotFoundError("No subscription package found.")
	}
	return subPackages, nil
}
