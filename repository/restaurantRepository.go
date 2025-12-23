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

func (rr *RestaurantRepository) CreateRestaurant(ctx context.Context, name string, location string, description *string, email *string, phone *string, postalCode string, rType enum.RestaurantType, notes string, subPackageId int64, images []string, tenantId int64) (*model.Restaurant, error) {
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
	if err = tx.QueryRowContext(ctx, query, name, location, description, email, phone, postalCode, rType, ar, true, notes, subPackageId, tenantId).Scan(
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

func (rr *RestaurantRepository) GetRestaurantByName(ctx context.Context, name string) (*model.Restaurant, error) {
	var err error
	var r model.Restaurant
	query := "SELECT * FROM restaurants WHERE name = $1 LIMIT 1"
	err = rr.Db.QueryRowContext(ctx, query, name).Scan(
		&r.Id,
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No restaurant found.")
		}	
		return nil, err
	} 
	return &r, nil
}

func (rr *RestaurantRepository) GetRestaurantImages(ctx context.Context, restaurantId int64) ([]model.RestaurantImage, error) {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.QueryContext(ctx, "SELECT * FROM restaurant_images WHERE restaurant_id = $1", restaurantId)
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

func (rr *RestaurantRepository) IsRestaurantNameExist(ctx context.Context, name string) (bool, error) {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.QueryContext(ctx, "SELECT id FROM restaurants WHERE name = $1 LIMIT 1", name)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}


func (rr *RestaurantRepository) IsRestaurantExist(ctx context.Context, restaurantId int64) (bool, error) {
	var err error
	var rows *sql.Rows
	rows, err = rr.Db.QueryContext(ctx, "SELECT id FROM restaurants WHERE id = $1 LIMIT 1", restaurantId)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (rr *RestaurantRepository) GetRestaurantById(ctx context.Context, id int64) (*model.Restaurant, error) {
	var err error
	var r model.Restaurant
	query := "SELECT * FROM restaurants WHERE id = $1 LIMIT 1"
	err = rr.Db.QueryRowContext(ctx, query, id).Scan(
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No restaurant found.")
		}
		return nil, err
	}
	return &r, nil
}

func (rr *RestaurantRepository) GetRestaurantsByTenantId(ctx context.Context, tenantId int64) ([]model.Restaurant, error) {
	var err error
	var rows *sql.Rows
	query := "SELECT * FROM restaurants WHERE tenant_id = $1"
	rows, err = rr.Db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return nil, err
	}
	var restaurants []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		if err = rows.Scan(
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
		restaurants = append(restaurants, r)
	}
	if len(restaurants) == 0 {
		return nil, appErr.NotFoundError("No restaruant found.")
	}
	return restaurants, nil
}

func (rr *RestaurantRepository) ListAllRestaurants(ctx context.Context) ([]model.Restaurant, error) {
	var err error
	var rows *sql.Rows
	query := "SELECT * FROM restaurants ORDER BY id ASC"
	rows, err = rr.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var restaurants []model.Restaurant
	for rows.Next() {
		var r model.Restaurant
		if err = rows.Scan(
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
		restaurants = append(restaurants, r)
	}
	if len(restaurants) == 0 {
		return nil, appErr.NotFoundError("No restaurant found.")
	}
	return restaurants, nil
}


func (rr *RestaurantRepository) GetRestaurantImageById(ctx context.Context, id int64) (*model.RestaurantImage, error) {
	var err error
	var i model.RestaurantImage
	query := "SELECT * FROM restaurant_images WHERE id = $1 LIMIT 1"
	err = rr.Db.QueryRowContext(ctx, query, id).Scan(
		&i.Id, 
		&i.Image,
		&i.CreatedAt,
		&i.RestaurantId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No image found.")
		}
		return nil, err
	}
	return &i, nil
}


func (rr *RestaurantRepository) ListAllRestaurantImages(ctx context.Context) ([]model.RestaurantImage, error) {
	var err error
	var rows *sql.Rows
	query := "SELECT * FROM restaurant_images ORDER BY id ASC"
	rows, err = rr.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var images []model.RestaurantImage
	for rows.Next() {
		var i model.RestaurantImage
		if err = rows.Scan(
			&i.Id, 
			&i.Image, 
			&i.CreatedAt, 
			&i.RestaurantId,
		); err != nil {
			return nil, err
		}
		images = append(images, i)
	}
	if len(images) == 0 {
		return nil, appErr.NotFoundError("No restaurant image found.")
	}
	return images, nil
}


func (rr *RestaurantRepository) AddNewRestaurantItem(ctx context.Context, name string, description *string, price decimal.Decimal, itemType enum.ItemType, image *string, notes *string, restaurantId int64) (*model.RestaurantItem, error) {
	var err error
	var tx *sql.Tx
	tx, err = rr.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var i model.RestaurantItem
	query := "INSERT INTO restaurant_items (name, description, price, type, status, image, notes, restaurant_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *"
	if err = tx.QueryRowContext(ctx, query, name, description, price, itemType, enum.ITEM_AVAILABLE, image, notes, restaurantId).Scan(
		&i.Id,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Type,
		&i.Status,
		&i.Image,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RestaurantId,
	); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &i, nil
}

