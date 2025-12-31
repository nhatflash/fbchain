package repository

import (
	"context"
	"database/sql"
	"github.com/nhatflash/fbchain/enum"
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


func (ror *RestaurantOrderRepository) CreateInitialRestaurantOrder(ctx context.Context, rOrder *model.RestaurantOrder, rOItems []model.RestaurantOrderItem) (*model.RestaurantOrder, error) {
	var err error
	var tx *sql.Tx
	tx, err = ror.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	defer tx.Rollback()

	query := "INSERT INTO restaurant_orders (restaurant_id, table_id, amount, status, notes) VALUES ($1, $2, $3, $4, $5) RETURNING *"
	var o model.RestaurantOrder
	if err = tx.QueryRowContext(ctx, query, rOrder.RestaurantId, rOrder.TableId, rOrder.Amount, enum.R_ORDER_PENDING, rOrder.Notes).Scan(
		&o.Id,
		&o.RestaurantId,
		&o.TableId,
		&o.Amount,
		&o.Status,
		&o.Notes,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return nil, err
	}

	var items []model.RestaurantOrderItem
	for _, item := range rOItems {
		var i model.RestaurantOrderItem
		itemQuery := "INSERT INTO restaurant_order_items (restaurant_order_id, item_id, quantity, total) VALUES ($1, $2, $3, $4) RETURNING *"
		if err = tx.QueryRowContext(ctx, itemQuery, o.Id, item.ItemId, item.Quantity, item.Total).Scan(
			&i.Id,
			&i.ROrderId,
			&i.ItemId,
			&i.Quantity,
			&i.Total,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	o.Items = items

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &o, nil
}