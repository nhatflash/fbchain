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

func (or *OrderRepository) CreateInitialOrder(rId int64, sId int64, amount *decimal.Decimal, tId int64) (error) {
	var err error
	ctx := context.Background()
	var tx *sql.Tx
	tx, err = or.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "INSERT INTO orders (status, amount, tenant_id, restaurant_id, subscription_id) VALUES ($1, $2, $3, $4, $5)", enum.ORDER_PENDING, amount, tId, rId, sId)
	if err != nil {
		return err
	}
	return tx.Commit()
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
		err = rows.Scan(&o.Id, &o.OrderDate, &o.Status, &o.Amount, &o.UpdatedAt, &o.TenantId, &o.RestaurantId, &o.SubPackageId)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if len(orders) == 0 {
		return nil, appErr.NotFoundError("No order found.")
	}
	return &orders[0], nil
}
