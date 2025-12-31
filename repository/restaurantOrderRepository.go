package repository

import (
	"context"
	"database/sql"

	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/model"
)


type RestaurantOrderRepository struct {
	Db  			*sql.DB
}


func NewRestaurantOrderRepositoty(db *sql.DB) *RestaurantOrderRepository {
	return &RestaurantOrderRepository{
		Db : db,
	}
}


func (ror *RestaurantOrderRepository) CreateInitialRestaurantOrder(ctx context.Context, tableId int64, req *client.CreateRestaurantOrderRequest) (*model.RestaurantOrder, error) {
	var err error
	var tx *sql.Tx
	tx, err = ror.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	defer tx.Rollback()

	query := "INSERT INTO restaurant_orders (id,restaurant_id, table_id, amount, status, notes) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *"
	return nil, nil
}