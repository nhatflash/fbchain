package repository

import (
	"database/sql"

	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
	"github.com/shopspring/decimal"
)

func CheckSubPackageNameExists(name string, db *sql.DB) bool {
	rows, rowErr := db.Query("SELECT name FROM sub_packages WHERE name = $1", name)
	if rowErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func CreateSubPackage(name string, description string, durationMonth int, price decimal.Decimal, image string, db *sql.DB) (*model.SubPackage, error) {
	var err error
	_, err = db.Exec("INSERT INTO sub_packages (name, description, duration_month, price, is_active, image) VALUES ($1, $2, $3, $4, $5, $6)", name, description, durationMonth, price, true, image)

	if err != nil {
		return nil, err
	}

	var s *model.SubPackage
	s, err = GetSubPackageByName(name, db)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func GetSubPackageByName(name string, db *sql.DB) (*model.SubPackage, error) {
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM sub_packages WHERE name = $1", name)
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

func AnySubPackageExists(db *sql.DB) (bool, error) {
	rows, rowErr := db.Query("SELECT id FROM sub_packages")
	if rowErr != nil {
		return false, rowErr
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func IsSubPackageExist(sId int64, db *sql.DB) bool {
	rows, rowErr := db.Query("SELECT id FROM sub_packages WHERE id = $1", sId)
	if rowErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func GetSubPackageById(sId int64, db *sql.DB) (*model.SubPackage, error) {
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM sub_packages WHERE id = $1", sId)
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
