package service

import (
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/enum"
)

type IPaymentService interface {
	HandleCashPayment(orderId int64, notes *string) error
	HandleVnPayPayment(orderId int64, status enum.PaymentStatus, bankCode *string, notes *string) error
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

func (ps *PaymentService) HandleCashPayment(orderId int64, notes *string) error {
	var err error
	var o *model.Order
	o, err = ps.OrderRepo.GetOrderById(orderId)
	if err != nil {
		return err
	}
	if err = ps.PaymentRepo.CreateCashPayment(orderId, o.Amount, notes); err != nil {
		return err
	}
	return nil
}

func (ps *PaymentService) HandleVnPayPayment(orderId int64, status enum.PaymentStatus, bankCode *string, notes *string) error {
	var err error
	var o *model.Order
	o, err = ps.OrderRepo.GetOrderById(orderId)
	if err != nil {
		return err
	}
	if err = ps.PaymentRepo.CreateOnlinePayment(orderId, o.Amount, enum.PAYMENT_VNPAY, status, bankCode, notes); err != nil {
		return err
	}
	return nil
}
