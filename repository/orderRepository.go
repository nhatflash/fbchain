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

	var o model.Order
	query := "INSERT INTO orders (status, amount, tenant_id, restaurant_id, subscription_id) VALUES ($1, $2, $3, $4, $5) RETURNING *"

	if err = tx.QueryRowContext(ctx, query, enum.ORDER_PENDING, amount, tId, restaurantId, subPackageId).Scan(
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
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &o, nil
}

func (or *OrderRepository) FindLatestTenantOrder(ctx context.Context, tenantId int64) (*model.Order, error) {
	query := "SELECT * FROM orders WHERE tenant_id = $1 ORDER BY order_date DESC LIMIT 1"
	var o model.Order
	err := or.Db.QueryRowContext(ctx, query, tenantId).Scan(
		&o.Id,
		&o.OrderDate,
		&o.Status,
		&o.Amount,
		&o.UpdatedAt,
		&o.TenantId,
		&o.RestaurantId,
		&o.SubPackageId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No order found.")
		}
		return nil, err
	}
	return &o, nil
}


func (or *OrderRepository) FindOrderById(ctx context.Context, id int64) (*model.Order, error) {
	query := "SELECT * FROM orders WHERE id = $1 LIMIT 1"
	var o model.Order
	err := or.Db.QueryRowContext(ctx, query, id).Scan(
		&o.Id,
		&o.OrderDate,
		&o.Status,
		&o.Amount,
		&o.UpdatedAt,
		&o.TenantId,
		&o.RestaurantId,
		&o.SubPackageId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No order found.")
		}
		return nil, err
	}
	return &o, nil
}


func (or *OrderRepository) FindAllOrders(ctx context.Context) ([]model.Order, error) {
	var err error
	var rows *sql.Rows
	query := "SELECT * FROM orders"
	rows, err = or.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err = rows.Scan(
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
		orders = append(orders, o)
	}
	if len(orders) == 0 {
		return nil, appErr.NotFoundError("No order found.")
	}
	return orders, nil
}


func (or *OrderRepository) FindOrdersByTenantId(ctx context.Context, tenantId int64) ([]model.Order, error) {
	var err error
	var rows *sql.Rows
	query := "SELECT * FROM orders WHERE tenant_id = $1"
	rows, err = or.Db.QueryContext(ctx, query, tenantId)
	if err != nil {
		return nil, err
	}
	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err = rows.Scan(
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
		orders = append(orders, o)
	}
	if len(orders) == 0 {
		return nil, appErr.NotFoundError("No order found.")
	}
	return orders, nil
}


func (or *OrderRepository) FinishOrder(ctx context.Context, orderId int64) (*model.Order, error) {
	var err error
	var tx *sql.Tx
	tx, err = or.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var o model.Order
	query := "UPDATE orders SET status = $1 WHERE id = $2"
	if err = tx.QueryRowContext(ctx, query, enum.ORDER_COMPLETED, orderId).Scan(
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
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &o, nil
}
