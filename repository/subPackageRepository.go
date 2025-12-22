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

func (spr *SubPackageRepository) CheckSubPackageNameExists(name string) bool {
	rows, rowErr := spr.Db.Query("SELECT name FROM sub_packages WHERE name = $1", name)
	if rowErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func (spr *SubPackageRepository) CreateSubPackage(ctx context.Context, name string, description *string, durationMonth int, price decimal.Decimal, image *string) (*model.SubPackage, error) {
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
		&s.Image,
	); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &s, nil
}

func (spr *SubPackageRepository) GetSubPackageByName(name string) (*model.SubPackage, error) {
	var err error
	var rows *sql.Rows
	rows, err = spr.Db.Query("SELECT * FROM sub_packages WHERE name = $1", name)
	if err != nil {
		return nil, err
	}
	var subPackages []model.SubPackage
	for rows.Next() {
		var s model.SubPackage
		err = rows.Scan(&s.Id, &s.Name, &s.Description, &s.DurationMonth, &s.Price, &s.IsActive, &s.Image)
		if err != nil {
			return nil, err
		}
		subPackages = append(subPackages, s)
	}
	if len(subPackages) == 0 {
		return nil, appErr.NotFoundError("No subscription package found.")
	}
	return &subPackages[0], nil
}

func (spr *SubPackageRepository) AnySubPackageExists() (bool, error) {
	rows, rowErr := spr.Db.Query("SELECT id FROM sub_packages")
	if rowErr != nil {
		return false, rowErr
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (spr *SubPackageRepository) IsSubPackageExist(sId int64) bool {
	rows, rowErr := spr.Db.Query("SELECT id FROM sub_packages WHERE id = $1", sId)
	if rowErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func (spr *SubPackageRepository) GetSubPackageById(sId int64) (*model.SubPackage, error) {
	var err error
	var rows *sql.Rows
	rows, err = spr.Db.Query("SELECT * FROM sub_packages WHERE id = $1", sId)
	if err != nil {
		return nil, err
	}
	var subPackages []model.SubPackage
	for rows.Next() {
		var s model.SubPackage
		err = rows.Scan(&s.Id, &s.Name, &s.Description, &s.DurationMonth, &s.Price, &s.IsActive, &s.Image)
		if err != nil {
			return nil, err
		}
		subPackages = append(subPackages, s)
	}
	if len(subPackages) == 0 {
		return nil, appErr.NotFoundError("No subscription found.")
	}
	return &subPackages[0], nil
}
