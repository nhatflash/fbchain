package repository

import (
	"context"
	"database/sql"

	"github.com/nhatflash/fbchain/enum"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
	"github.com/shopspring/decimal"
)

type OrderRepository struct {
	Db 			*sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		Db: db,
	}
}

func (or *OrderRepository) CreateInitialOrder(ctx context.Context, restaurantId int64, subPackageId int64, amount *decimal.Decimal, tId int64) (*model.Order, error) {
	var err error
	var tx *sql.Tx
	tx, err = or.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var newOrder model.Order
	query := "INSERT INTO orders (status, amount, tenant_id, restaurant_id, subscription_id) VALUES ($1, $2, $3, $4, $5) RETURNING *"

	if err = tx.QueryRowContext(ctx, query, enum.ORDER_PENDING, amount, tId, restaurantId, subPackageId).Scan(
		&newOrder.Id,
		&newOrder.OrderDate,
		&newOrder.Status,
		&newOrder.Amount,
		&newOrder.UpdatedAt,
		&newOrder.TenantId,
		&newOrder.RestaurantId,
		&newOrder.SubPackageId,
	); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &newOrder, nil
}

func (or *OrderRepository) GetLatestTenantOrder(ctx context.Context, tenantId int64) (*model.Order, error) {
	var err error

	query := "SELECT * FROM orders WHERE tenant_id = $1 ORDER BY order_date DESC LIMIT 1"
	var order model.Order
	if err = or.Db.QueryRowContext(ctx, query, tenantId).Scan(
		&order.Id,
		&order.OrderDate,
		&order.Status,
		&order.Amount,
		&order.UpdatedAt,
		&order.TenantId,
		&order.RestaurantId,
		&order.SubPackageId,
	); err != nil {
		return nil, err
	}
	if err != nil && err == sql.ErrNoRows {
		return nil, appErr.NotFoundError("No order found.")
	}
	return &order, nil
}


func (or *OrderRepository) GetOrderById(ctx context.Context, id int64) (*model.Order, error) {
	var err error
	var o model.Order
	query := "SELECT * FROM orders WHERE id = $1 LIMIT 1"
	if err = or.Db.QueryRowContext(ctx, query, id).Scan(
		&o.Id,
		&o.OrderDate,
		&o.Status,
		&o.Amount,
		&o.UpdatedAt,
		&o.TenantId,
		&o.RestaurantId,
		&o.SubPackageId,
	); err != nil {
		return nil, err
	}
	if err != nil && err == sql.ErrNoRows {
		return nil, appErr.NotFoundError("No order found.")
	}
	return &o, nil
}
