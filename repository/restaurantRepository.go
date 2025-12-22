package repository

import (
	"context"
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

func (rr *RestaurantRepository) CreateRestaurant(ctx context.Context, name string, location string, description *string, email *string, phone *string, postalCode string, rType *enum.RestaurantType, notes string, images []string, tenantId int64) (*model.Restaurant, error) {
	var tx *sql.Tx	
	var err error
	tx, err = rr.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var ar decimal.Decimal
	ar, err = decimal.NewFromString("0.0")
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO restaurants (name, location, description, contact_email, contact_phone, postal_code, type, avg_rating, is_active, notes, subscription_id, tenant_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *"
	var r model.Restaurant
	if err = tx.QueryRowContext(ctx, query, name, location, description, email, phone, postalCode, rType, ar, true, notes, 1, tenantId).Scan(
		&r.Id,
		&r.Name,
		&r.Location,
		&r.Description,
		&r.ContactEmail,
		&r.ContactPhone,
		&r.PostalCode,
		&r.Type,
		&r.AvgRating,
		&r.IsActive,
		&r.Notes,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.SubPackageId,
		&r.TenantId,
	); err != nil {
		return nil, err
	}
	if len(images) > 0 {
		imgQuery := "INSERT INTO restaurant_images (image, restaurant_id) VALUES ($1, $2)"
		for i := range images {
			_, err = tx.ExecContext(ctx, imgQuery, images[i], r.Id)
			if err != nil {
				return nil, err
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &r, nil
}

func (rr *RestaurantRepository) CreateRestaurantImages(ctx context.Context, rId int64, images []string) ([]model.RestaurantImage, error) {
	var err error
	var tx *sql.Tx
	tx, err = rr.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	
	var imgs []model.RestaurantImage
	query := "INSERT INTO restaurant_images (image, restaurant_id) VALUES ($1, $2) RETURNING *"
	for i := range images {
		var img model.RestaurantImage
		if err = tx.QueryRowContext(ctx, query, images[i], rId).Scan(
			&img.Id,
			&img.Image,
			&img.CreatedAt,
			&img.RestaurantId,
		); err != nil {
			return nil, err
		}
		imgs = append(imgs, img)
	}
	return imgs, nil
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

