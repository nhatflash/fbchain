package repository

import (
	"context"
	"database/sql"

	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
)

type StaffRepository struct {
	Db  	*sql.DB
}

func NewStaffRepository(db *sql.DB) *StaffRepository {
	return &StaffRepository{
		Db: db,
	}
}

func (sr *StaffRepository) FindStaffInfoById(ctx context.Context, id int64) (*model.Staff, error) {
	var err error
	var s model.Staff
	query := "SELECT * FROM staffs WHERE id = $1 LIMIT 1"
	err = sr.Db.QueryRowContext(ctx, query, id).Scan(
		&s.Id,
		&s.UserId,
		&s.Code,
		&s.ShiftType,
		&s.ShiftStart,
		&s.ShiftEnd,
		&s.Salary,
		&s.Notes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No staff info found.")
		}
		return nil, err
	}
	return &s, nil

}


func (sr *StaffRepository) FindStaffInfoByUserId(ctx context.Context, userId int64) (*model.Staff, error) {
	var err error
	var s model.Staff
	query := "SELECT * FROM staffs WHERE user_id = $1 LIMIT 1"
	err = sr.Db.QueryRowContext(ctx, query, userId).Scan(
		&s.Id, 
		&s.UserId,
		&s.Code,
		&s.ShiftType,
		&s.ShiftStart,
		&s.ShiftEnd,
		&s.Salary,
		&s.Notes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No staff info found.")
		}
		return nil, err
	}
	return &s, nil
}