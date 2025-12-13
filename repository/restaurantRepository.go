package repository

import (
	"database/sql"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/enum"
	"github.com/shopspring/decimal"
	appErr "github.com/nhatflash/fbchain/error"
)

func CreateRestaurant(name string, location string, description string, email string, phone string, postalCode string, rType *enum.RestaurantType, notes string, images []string, tenantId int64, db *sql.DB) (*model.Restaurant, error) {
	avgRating, dErr := decimal.NewFromString("0.0")
	if dErr != nil {
		return nil, dErr
	}
	_, insertErr := db.Exec("INSERT INTO restaurants (name, location, description, contact_email, contact_phone, postal_code, type, avg_rating, is_active, notes, subscription_id, tenant_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", name, location, description, email, phone, postalCode, rType, avgRating, true, notes, 1, tenantId)
	if insertErr != nil {
		return nil, insertErr
	}
	r, rErr := GetRestaurantByName(name, db)
	if rErr != nil {
		return nil, rErr
	}
	imageErr := CreateRestaurantImages(r.Id, images, db)
	if imageErr != nil {
		return nil, imageErr
	}
	rImages, rIErr := GetRestaurantImages(r.Id, db)
	if rIErr != nil {
		return nil, rIErr
	}
	r.Images = *rImages
	return r, nil
}

func CreateRestaurantImages(rId int64, images []string, db *sql.DB) error {
	for i := range images {
		_, insertErr := db.Exec("INSERT INTO restaurant_images (image, restaurant_id) VALUES ($1, $2)", images[i], rId)
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}


func GetRestaurantByName(name string, db *sql.DB) (*model.Restaurant, error) {
	rows, rowErr := db.Query("SELECT * FROM restaurants WHERE name = $1 LIMIT 1 ", name)
	if rowErr != nil {
		return nil, rowErr
	}
	var restaurants []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		scanErr := rows.Scan(&r.Id, &r.Name, &r.Location, &r.Description, &r.ContactEmail, &r.ContactPhone, &r.PostalCode, &r.Type, &r.AvgRating, &r.IsActive, &r.Notes, &r.CreatedAt, &r.UpdatedAt, &r.SubscriptionId, &r.TenantId)
		if scanErr != nil {
			return nil, scanErr
		}
		restaurants = append(restaurants, r)
	}
	if len(restaurants) == 0 {
		return nil, appErr.NotFoundError("No restaurant found.")
	}
	return &restaurants[0], nil
}


func GetRestaurantImages(rId int64, db *sql.DB) (*[]model.RestaurantImage, error) {
	rows, rowErr := db.Query("SELECT * FROM restaurant_images WHERE restaurant_id = $1", rId)
	if rowErr != nil {
		return nil, rowErr
	}
	var images []model.RestaurantImage
	for rows.Next() {
		var i model.RestaurantImage
		scanErr := rows.Scan(&i.Id, &i.Image, &i.CreatedAt, &i.RestaurantId)
		if scanErr != nil {
			return nil, scanErr
		}
		images = append(images, i)
	}
	if len(images) == 0 {
		return nil, appErr.NotFoundError("No image found.")
	}
	return &images, nil
}

func IsRestaurantNameExist(name string, db *sql.DB) bool {
	rows, rowErr := db.Query("SELECT id FROM restaurants WHERE name = $1 LIMIT 1", name)
	if rowErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}