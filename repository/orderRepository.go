package repository

import (
	"context"
	"database/sql"

	"github.com/nhatflash/fbchain/enum"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
	"github.com/shopspring/decimal"
)

func CreateInitialOrder(rId int64, sId int64, amount *decimal.Decimal, tId int64, db *sql.DB) (err error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	_, err = tx.ExecContext(ctx, "INSERT INTO orders (status, amount, tenant_id, restaurant_id, subscription_id) VALUES ($1, $2, $3, $4, $5)", enum.ORDER_PENDING, amount, tId, rId, sId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func GetLatestTenantOrder(tId int64, db *sql.DB) (*model.Order, error) {
	var rows *sql.Rows
	var err error
	rows, err = db.Query("SELECT * FROM orders WHERE tenant_id = $1 ORDER BY order_date", tId)
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
