package repository

import (
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

func (spr *SubPackageRepository) CreateSubPackage(name string, description string, durationMonth int, price decimal.Decimal, image string) (*model.SubPackage, error) {
	var err error
	_, err = spr.Db.Exec("INSERT INTO sub_packages (name, description, duration_month, price, is_active, image) VALUES ($1, $2, $3, $4, $5, $6)", name, description, durationMonth, price, true, image)

	if err != nil {
		return nil, err
	}

	var s *model.SubPackage
	s, err = spr.GetSubPackageByName(name)
	if err != nil {
		return nil, err
	}
	return s, nil
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
