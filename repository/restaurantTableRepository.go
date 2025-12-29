package repository

import (
	"context"
	"database/sql"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
)


type RestaurantTableRepository struct {
	Db 			*sql.DB
}


func NewRestaurantTableRepository(db *sql.DB) *RestaurantTableRepository {
	return &RestaurantTableRepository{
		Db: db,
	}
}


func (rtr *RestaurantTableRepository) AddNewRestaurantTable(ctx context.Context, restaurantId int64, label string, notes *string) (*model.RestaurantTable, error) {
	var err error
	var tx *sql.Tx
	tx, err = rtr.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var t model.RestaurantTable
	query := "INSERT INTO restaurant_tables (restaurant_id, notes) VALUES ($1, $2) RETURNING *"
	if err = tx.QueryRowContext(ctx, query, restaurantId, notes).Scan(
		&t.Id,
		&t.RestaurantId,
		&t.Label,
		&t.IsActive,
		&t.Notes,
		&t.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &t, nil
}


func (rtr *RestaurantTableRepository) FindAllRestaurantTables(ctx context.Context) ([]model.RestaurantTable, error) {
	var err error
	var rows *sql.Rows
	query := "SELECT * FROM restaurant_tables"
	rows, err = rtr.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var tables []model.RestaurantTable
	for rows.Next() {
		var t model.RestaurantTable
		if err = rows.Scan(
			&t.Id,
			&t.RestaurantId,
			&t.Label,
			&t.IsActive,
			&t.Notes,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	if len(tables) == 0 {
		return nil, appErr.NotFoundError("No restaurant tables found.")
	}
	return tables, nil
}

func (rtr *RestaurantTableRepository) FindRestaurantTableById(ctx context.Context, id int64) (*model.RestaurantTable, error) {
	var t model.RestaurantTable
	query := "SELECT * FROM restaurant_tables WHERE id = $1 LIMIT 1"
	err := rtr.Db.QueryRowContext(ctx, query, id).Scan(
		&t.Id,
		&t.RestaurantId,
		&t.Label,
		&t.IsActive,
		&t.Notes,
		&t.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No restaurant table found.")
		}
		return nil, err
	}
	return &t, nil
}


func (rtr *RestaurantTableRepository) CountRestaurantTableByRestaurantId(ctx context.Context, restaurantId int64) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM restaurant_tables WHERE restaurant_id = $1"
	if err := rtr.Db.QueryRowContext(ctx, query, restaurantId).Scan(
		&count,
	); err != nil {
		return 0, err
	}
	return count, nil
}


func (rtr *RestaurantTableRepository) FindRestaurantTablesByRestaurantId(ctx context.Context, restaurantId int64) ([]model.RestaurantTable, error) {
	var err error
	var rows *sql.Rows
	query := "SELECT * FROM restaurant_tables WHERE restaurant_id = $1"
	rows, err = rtr.Db.QueryContext(ctx, query, restaurantId)
	if err != nil {
		return nil, err
	}
	var tables []model.RestaurantTable
	for rows.Next() {
		var t model.RestaurantTable
		if err = rows.Scan(
			&t.Id,
			&t.RestaurantId,
			&t.Label,
			&t.IsActive,
			&t.Notes,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	if len(tables) == 0 {
		return nil, appErr.NotFoundError("No restaurant table found.")
	}
	return tables, nil
}