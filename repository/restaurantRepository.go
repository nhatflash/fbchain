package repository

import (
	"database/sql"

	"github.com/nhatflash/fbchain/enum"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
	"github.com/shopspring/decimal"
)

type RestaurantRepository struct {
	Db 			*sql.DB
}

func NewRestaurantRepository(db *sql.DB) *RestaurantRepository {
	return &RestaurantRepository{
		Db: db,
	}
}

func (rr *RestaurantRepository) CreateRestaurant(name string, location string, description string, email string, phone string, postalCode string, rType *enum.RestaurantType, notes string, images []string, tenantId int64) (*model.Restaurant, error) {
	var ar decimal.Decimal
	var err error

	ar, err = decimal.NewFromString("0.0")
	if err != nil {
		return nil, err
	}
	_, err = rr.Db.Exec("INSERT INTO restaurants (name, location, description, contact_email, contact_phone, postal_code, type, avg_rating, is_active, notes, subscription_id, tenant_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", name, location, description, email, phone, postalCode, rType, ar, true, notes, 1, tenantId)
	if err != nil {
		return nil, err
	}
	r, err := rr.GetRestaurantByName(name)
	if err != nil {
		return nil, err
	}
	err = rr.CreateRestaurantImages(r.Id, images)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rr *RestaurantRepository) CreateRestaurantImages(rId int64, images []string) error {
	var err error
	for i := range images {
		_, err = rr.Db.Exec("INSERT INTO restaurant_images (image, restaurant_id) VALUES ($1, $2)", images[i], rId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rr *RestaurantRepository) GetRestaurantByName(name string) (*model.Restaurant, error) {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.Query("SELECT * FROM restaurants WHERE name = $1 LIMIT 1 ", name)
	if err != nil {
		return nil, err
	}
	var restaurants []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		err = rows.Scan(&r.Id, &r.Name, &r.Location, &r.Description, &r.ContactEmail, &r.ContactPhone, &r.PostalCode, &r.Type, &r.AvgRating, &r.IsActive, &r.Notes, &r.CreatedAt, &r.UpdatedAt, &r.SubPackageId, &r.TenantId)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}
	if len(restaurants) == 0 {
		return nil, appErr.NotFoundError("No restaurant found.")
	}
	return &restaurants[0], nil
}

func (rr *RestaurantRepository) GetRestaurantImages(rId int64) ([]model.RestaurantImage, error) {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.Query("SELECT * FROM restaurant_images WHERE restaurant_id = $1", rId)
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
	return images, nil
}

func (rr *RestaurantRepository) IsRestaurantNameExist(name string) bool {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.Query("SELECT id FROM restaurants WHERE name = $1 LIMIT 1", name)
	if err != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}


func (rr *RestaurantRepository) IsRestaurantExist(rId int64) bool {
	rows, rowErr := rr.Db.Query("SELECT id FROM restaurants WHERE id = $1 LIMIT 1", rId)
	if rowErr != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func (rr *RestaurantRepository) GetRestaurantById(rId int64) (*model.Restaurant, error) {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.Query("SELECT * FROM restaurants WHERE id = $1 LIMIT 1", rId)
	if err != nil {
		return nil, err
	}
	var restaurants []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		err = rows.Scan(&r.Id, &r.Name, &r.Location, &r.Description, &r.ContactEmail, &r.ContactPhone, &r.PostalCode, &r.Type, &r.AvgRating, &r.IsActive, &r.Notes, &r.CreatedAt, &r.UpdatedAt, &r.SubPackageId, &r.TenantId)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}
	if len(restaurants) == 0 {
		return nil, appErr.NotFoundError("No restaurant found.")
	}
	return &restaurants[0], nil
}

func (rr *RestaurantRepository) GetRestaurantsByTenantId(tId int64) ([]model.Restaurant, error) {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.Query("SELECT * FROM restaurants WHERE tenant_id = $1", tId)
	if err != nil {
		return nil, err
	}
	var restaurants []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		err = rows.Scan(&r.Id, &r.Name, &r.Location, &r.Description, &r.ContactEmail, &r.ContactPhone, &r.PostalCode, &r.Type, &r.AvgRating, &r.IsActive, &r.Notes, &r.CreatedAt, &r.UpdatedAt, &r.SubPackageId, &r.TenantId)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}
	if len(restaurants) == 0 {
		return nil, appErr.NotFoundError("No restaruant found.")
	}
	return restaurants, nil
}

func (rr *RestaurantRepository) ListAllRestaurants() ([]model.Restaurant, error) {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.Query("SELECT * FROM restaurants ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	var restaurants []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		err = rows.Scan(&r.Id, &r.Name, &r.Location, &r.Description, &r.ContactEmail, &r.ContactPhone, &r.PostalCode, &r.Type, &r.AvgRating, &r.IsActive, &r.Notes, &r.CreatedAt, &r.UpdatedAt, &r.SubPackageId, &r.TenantId)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}
	if len(restaurants) == 0 {
		return nil, appErr.NotFoundError("No restaurant found.")
	}
	return restaurants, nil
}


func (rr *RestaurantRepository) GetRestaurantImageById(id int64) (*model.RestaurantImage, error) {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.Query("SELECT * FROM restaurant_images WHERE id = $1 LIMIT 1", id)
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
		return nil, appErr.NotFoundError("No restaurant image found.")
	}
	return &images[0], nil
}


func (rr *RestaurantRepository) ListAllRestaurantImages() ([]model.RestaurantImage, error) {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.Query("SELECT * FROM restaurant_images ORDER BY id ASC")
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
		return nil, appErr.NotFoundError("No restaurant image found.")
	}
	return images, nil
}

