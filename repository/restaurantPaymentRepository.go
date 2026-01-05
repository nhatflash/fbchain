package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/nhatflash/fbchain/enum"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
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


func (rpr *RestaurantPaymentRepository) ConfirmCashedPayment(ctx context.Context, orderId int64, paymentId int64) error {
	var err error
	var tx *sql.Tx
	tx, err = rpr.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	paymentQuery := "UPDATE restaurant_payments SET is_cashed = $1 WHERE id = $2"
	_, err = tx.ExecContext(ctx, paymentQuery, true, paymentId)
	if err != nil {
		return err
	}

	orderQuery := "UPDATE restaurant_orders SET status = $1, updated_at = $2 WHERE id = $3"
	_, err = tx.ExecContext(ctx, orderQuery, enum.R_ORDER_WORKING, time.Now(), orderId)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}


func (rpr *RestaurantPaymentRepository) FindSuccessfulPaymentOnOrderId(ctx context.Context, orderId int64) (*model.RestaurantPayment, error) {
	var err error
	var p model.RestaurantPayment
	query := "SELECT * FROM restaurant_payments WHERE order_id = $1 ORDER BY id DESC LIMIT 1"
	err = rpr.Db.QueryRowContext(ctx, query, orderId).Scan(
		&p.Id,
		&p.ROrderId,
		&p.Amount,
		&p.BankCode,
		&p.Method,
		&p.Status,
		&p.IsCashed,
		&p.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.NotFoundError("No restaurant payment found.")
		}
		return nil, err
	}
	return &p, nil
}