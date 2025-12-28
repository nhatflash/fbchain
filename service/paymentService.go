package service

import (
	"context"

	"github.com/nhatflash/fbchain/enum"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
)

type IPaymentService interface {
	HandleCashPayment(ctx context.Context, orderId int64, notes *string) error
	HandleVnPayPayment(ctx context.Context, orderId int64, status enum.PaymentStatus, bankCode *string, notes *string) error
}

type PaymentService struct {
	PaymentRepo 		*repository.PaymentRepository
	OrderRepo 			*repository.OrderRepository
}

func NewPaymentService(pr *repository.PaymentRepository, or *repository.OrderRepository) IPaymentService {
	return &PaymentService{
		PaymentRepo: pr,
		OrderRepo: or,
	}
}

func (ps *PaymentService) HandleCashPayment(ctx context.Context, orderId int64, notes *string) error {
	var err error
	var o *model.Order
	o, err = ps.OrderRepo.FindOrderById(ctx, orderId)
	if err != nil {
		return err
	}
	if err = ps.PaymentRepo.CreateCashPayment(ctx, orderId, o.Amount, notes); err != nil {
		return err
	}
	return nil
}

func (ps *PaymentService) HandleVnPayPayment(ctx context.Context, orderId int64, status enum.PaymentStatus, bankCode *string, notes *string) error {
	var err error
	var o *model.Order
	o, err = ps.OrderRepo.FindOrderById(ctx, orderId)
	if err != nil {
		return err
	}
	if err = ps.PaymentRepo.CreateOnlinePayment(ctx, orderId, o.Amount, enum.PAYMENT_VNPAY, status, bankCode, notes); err != nil {
		return err
	}
	return nil
}
