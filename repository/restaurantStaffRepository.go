package repository

import "database/sql"

type RestaurantStaffRepository struct {
	Db *sql.DB
}


func NewRestaurantStaffRepository(db *sql.DB) *RestaurantStaffRepository {
	return &RestaurantStaffRepository{
		Db: db,
	}
}


