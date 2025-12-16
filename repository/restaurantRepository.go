package repository

import (
	"database/sql"

	"github.com/nhatflash/fbchain/enum"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
	"github.com/shopspring/decimal"
)

func CreateRestaurant(name string, location string, description string, email string, phone string, postalCode string, rType *enum.RestaurantType, notes string, images []string, tenantId int64, db *sql.DB) (*model.Restaurant, error) {
	var ar decimal.Decimal
	var err error

	ar, err = decimal.NewFromString("0.0")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("INSERT INTO restaurants (name, location, description, contact_email, contact_phone, postal_code, type, avg_rating, is_active, notes, subscription_id, tenant_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", name, location, description, email, phone, postalCode, rType, ar, true, notes, 1, tenantId)
	if err != nil {
		return nil, err
	}
	r, err := GetRestaurantByName(name, db)
	if err != nil {
		return nil, err
	}
	err = CreateRestaurantImages(r.Id, images, db)
	if err != nil {
		return nil, err
	}

	var rImgs *[]model.RestaurantImage
	rImgs, err = GetRestaurantImages(r.Id, db)
	if err != nil {
		return nil, err
	}
	r.Images = *rImgs
	return r, nil
}

func CreateRestaurantImages(rId int64, images []string, db *sql.DB) error {
	var err error
	for i := range images {
		_, err = db.Exec("INSERT INTO restaurant_images (image, restaurant_id) VALUES ($1, $2)", images[i], rId)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetRestaurantByName(name string, db *sql.DB) (*model.Restaurant, error) {
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM restaurants WHERE name = $1 LIMIT 1 ", name)
	if err != nil {
		return nil, err
	}
	var restaurants []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		err = rows.Scan(&r.Id, &r.Name, &r.Location, &r.Description, &r.ContactEmail, &r.ContactPhone, &r.PostalCode, &r.Type, &r.AvgRating, &r.IsActive, &r.Notes, &r.CreatedAt, &r.UpdatedAt, &r.SubscriptionId, &r.TenantId)
		if err != nil {
			return nil, err
		}
		var rImgs *[]model.RestaurantImage
		rImgs, err = GetRestaurantImages(r.Id, db)
		if err != nil {
			return nil, err
		}
		r.Images = *rImgs
		restaurants = append(restaurants, r)
	}
	if len(restaurants) == 0 {
		return nil, appErr.NotFoundError("No restaurant found.")
	}
	return &restaurants[0], nil
}

func GetRestaurantImages(rId int64, db *sql.DB) (*[]model.RestaurantImage, error) {
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM restaurant_images WHERE restaurant_id = $1", rId)
	if err != nil {
		return nil, err
	}
	var images []model.RestaurantImage
	for rows.Next() {
		var i model.RestaurantImage
		err = rows.Scan(&i.Id, &i.Image, &i.CreatedAt, &i.RestaurantId)
		if err != nil {
			return nil, err
		}
		images = append(images, i)
	}
	if len(images) == 0 {
		return nil, appErr.NotFoundError("No image found.")
	}
	return &images, nil
}

func IsRestaurantNameExist(name string, db *sql.DB) bool {
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT id FROM restaurants WHERE name = $1 LIMIT 1", name)
	if err != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func IsRestaurantExist(rId int64, db *sql.DB) bool {
	rows, rowErr := db.Query("SELECT id FROM restaurants WHERE id = $1 LIMIT 1", rId)
	if rowErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func GetRestaurantById(rId int64, db *sql.DB) (*model.Restaurant, error) {
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM restaurants WHERE id = $1 LIMIT 1", rId)
	if err != nil {
		return nil, err
	}
	var restaurants []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		err = rows.Scan(&r.Id, &r.Name, &r.Location, &r.Description, &r.ContactEmail, &r.ContactPhone, &r.PostalCode, &r.Type, &r.AvgRating, &r.IsActive, &r.Notes, &r.CreatedAt, &r.UpdatedAt, &r.SubscriptionId, &r.TenantId)
		if err != nil {
			return nil, err
		}
		var rImgs *[]model.RestaurantImage
		rImgs, err = GetRestaurantImages(r.Id, db)
		if err != nil {
			return nil, err
		}
		r.Images = *rImgs
		restaurants = append(restaurants, r)
	}
	if len(restaurants) == 0 {
		return nil, appErr.NotFoundError("No restaurant found.")
	}
	return &restaurants[0], nil
}

func GetRestaurantsByTenantId(tId int64, db *sql.DB) ([]model.Restaurant, error) {
	var err error
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM restaurants WHERE tenant_id = $1", tId)
	if err != nil {
		return nil, err
	}
	var restaurants []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		err = rows.Scan(&r.Id, &r.Name, &r.Location, &r.Description, &r.ContactEmail, &r.ContactPhone, &r.PostalCode, &r.Type, &r.AvgRating, &r.IsActive, &r.Notes, &r.CreatedAt, &r.UpdatedAt, &r.SubscriptionId, &r.TenantId)
		if err != nil {
			return nil, err
		}
		var rImgs *[]model.RestaurantImage
		rImgs, err = GetRestaurantImages(r.Id, db)
		if err != nil {
			return nil, err
		}
		r.Images = *rImgs
		restaurants = append(restaurants, r)
	}
	if len(restaurants) == 0 {
		return nil, appErr.NotFoundError("No restaruant found.")
	}
	return restaurants, nil
}
