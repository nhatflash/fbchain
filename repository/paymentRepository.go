package repository

import (
	"context"
	"database/sql"
	"github.com/shopspring/decimal"
	"github.com/nhatflash/fbchain/enum"
)

type PaymentRepository struct {
	Db 			*sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{
		Db: db,
	}
}

func (pr *PaymentRepository) CreateOnlinePayment(ctx context.Context, orderId int64, amount decimal.Decimal, method enum.PaymentMethod, status enum.PaymentStatus, bankCode *string, notes *string) error {
	var err error
	var tx *sql.Tx
	tx, err = pr.Db.BeginTx(ctx, nil)

	defer tx.Rollback()
	query := "INSERT INTO payments (order_id, amount, method, bank_code, status, notes) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err = tx.ExecContext(ctx, query, orderId, amount, method, bankCode, status, notes)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (pr *PaymentRepository) CreateCashPayment(ctx context.Context, orderId int64, amount decimal.Decimal, notes *string) error {
	var err error
	var tx *sql.Tx
	tx, err = pr.Db.BeginTx(ctx, nil)

	defer tx.Rollback()
	query := "INSERT INTO payments (order_id, amount, method, bank_code, status, notes) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = tx.ExecContext(ctx, query, orderId, amount, enum.PAYMENT_CASH, nil, enum.PAYMENT_SUCCESS, notes)
	if err != nil {
		return err
	} 
	return tx.Commit()
} 
