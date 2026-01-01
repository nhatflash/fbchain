package repository

import (
	"context"
	"database/sql"
	"time"
	appErr "github.com/nhatflash/fbchain/error"
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


func (ror *RestaurantOrderRepository) DeleteExpiredPendingOrders(ctx context.Context) error {
	var err error
	var tx *sql.Tx
	tx, err = ror.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	
	expiration := time.Now().Add(-15 * time.Minute)
	query := "DELETE FROM restaurant_orders WHERE created_at < $1 AND status = $2"
	_, err = tx.ExecContext(ctx, query, expiration, enum.R_ORDER_PENDING)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}


func (ror *RestaurantOrderRepository) FindRestaurantOrderById(ctx context.Context, id int64) (*model.RestaurantOrder, error) {
	var err error
	var o model.RestaurantOrder
	query := "SELECT * FROM restaurant_orders WHERE id = $1 LIMIT 1" 
	err = ror.Db.QueryRowContext(ctx, query, id).Scan(
		&o.Id,
		&o.RestaurantId,
		&o.TableId,
		&o.Amount,
		&o.Notes,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No restaurant order found.")
		}
		return nil, err
	}
	var items []model.RestaurantOrderItem
	items, err = ror.FindRestaurantOrderItemsByOrderId(ctx, id)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return &o, nil
}


func (ror *RestaurantOrderRepository) FindRestaurantOrderItemsByOrderId(ctx context.Context, orderId int64) ([]model.RestaurantOrderItem, error) {
	var err error
	var rows *sql.Rows
	query := "SELECT * FROM restaurant_order_items WHERE restaurant_order_id = $1"
	rows, err = ror.Db.QueryContext(ctx, query, orderId)
	if err != nil {
		return nil, err
	}
	var items []model.RestaurantOrderItem
	for rows.Next() {
		var i model.RestaurantOrderItem
		if err = rows.Scan(
			&i.Id,
			&i.ROrderId,
			&i.ItemId,
			&i.Quantity,
			&i.Total,
		); err != nil {
			return nil, err
		}
	}
	if len(items) == 0 {
		return nil, appErr.NotFoundError("No restaurant order found.")
	}
	return items, nil
}