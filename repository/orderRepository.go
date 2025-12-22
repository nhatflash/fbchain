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

func (or *OrderRepository) GetLatestTenantOrder(tId int64) (*model.Order, error) {
	var rows *sql.Rows
	var err error
	rows, err = or.Db.Query("SELECT * FROM orders WHERE tenant_id = $1 ORDER BY order_date", tId)
	if err != nil {
		return nil, err
	}
	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err = rows.Scan(&o.Id, &o.OrderDate, &o.Status, &o.Amount, &o.UpdatedAt, &o.TenantId, &o.RestaurantId, &o.SubPackageId); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if len(orders) == 0 {
		return nil, appErr.NotFoundError("No order found.")
	}
	return &orders[0], nil
}


func (or *OrderRepository) GetOrderById(oId int64) (*model.Order, error) {
	var rows *sql.Rows
	var err error
	rows, err = or.Db.Query("SELECT * FROM orders WHERE id = $1 LIMIT 1", oId)
	if err != nil {
		return nil, err
	}
	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err = rows.Scan(&o.Id, &o.OrderDate, &o.Status, &o.Amount, &o.UpdatedAt, &o.TenantId, &o.RestaurantId, &o.SubPackageId); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if len(orders) == 0 {
		return nil, appErr.NotFoundError("No order found.")
	}
	return &orders[0], nil
}
