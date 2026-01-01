package repository

import (
	"context"
	"database/sql"

	"github.com/nhatflash/fbchain/enum"
	"github.com/shopspring/decimal"
)

type RestaurantPaymentRepository struct {
	Db 		 		*sql.DB	
}


func NewRestaurantPaymentRepository(db *sql.DB) *RestaurantPaymentRepository {
	return &RestaurantPaymentRepository{
		Db: db,
	}
}


func (rpr *RestaurantPaymentRepository) HandleCashPayment(ctx context.Context, orderId int64, amount decimal.Decimal) error {
	var err error
	var tx *sql.Tx
	tx, err = rpr.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := "INSERT INTO restaurant_payments (restaurant_order_id, amount, method, status) VALUES ($1, $2, $3, $4) RETURNING *"
	_, err = tx.ExecContext(ctx, query, orderId, amount, enum.PAYMENT_CASH, enum.PAYMENT_SUCCESS)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}