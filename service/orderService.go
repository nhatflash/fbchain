package service

import (
	"context"

	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
)

type IOrderService interface {
	HandlePaySubPackage(ctx context.Context, req *client.PaySubPackageRequest, tenantId int64) (*client.OrderResponse, error)
	FindOrderById(ctx context.Context, id int64) (*model.Order, error)
	FindAllOrders(ctx context.Context) ([]model.Order, error)
	FindOrdersByTenantId(ctx context.Context, tenantId int64) ([]model.Order, error)
}

type OrderService struct {
	RestaurantRepo 		*repository.RestaurantRepository
	SubPackageRepo 		*repository.SubPackageRepository
	OrderRepo 			*repository.OrderRepository
}

func NewOrderService(rr *repository.RestaurantRepository, 
					spr *repository.SubPackageRepository, 
					or *repository.OrderRepository) IOrderService {
	return &OrderService{
		RestaurantRepo: rr,
		SubPackageRepo: spr,
		OrderRepo: or,
	}
}

func (os *OrderService) HandlePaySubPackage(ctx context.Context, req *client.PaySubPackageRequest, tenantId int64) (*client.OrderResponse, error) {
	var err error
	var r *model.Restaurant
	var s *model.SubPackage

	r, s, err = checkRestaurantAndSubPackageExist(ctx, *req.RestaurantId, *req.SubPackageId, os.RestaurantRepo, os.SubPackageRepo)
	if err != nil {
		return nil, err
	}
	if err = checkIfRequestedRestaurantBelongToTenant(r, tenantId); err != nil {
		return nil, err
	}
	if isRestaurantSubPackageMatchTheRequestedPaySubPackage(r, s.Id) {
		return nil, appErr.BadRequestError("The requested subscription package is already registered on this restaurant.")
	}
	var newOrder *model.Order
	newOrder, err = os.OrderRepo.CreateInitialOrder(ctx, *req.RestaurantId, *req.SubPackageId, &s.Price, tenantId)
	if err != nil {
		return nil, err
	}
	return helper.MapToOrderResponse(newOrder), nil
}


func (os *OrderService) FindOrderById(ctx context.Context, id int64) (*model.Order, error) {
	o, err := os.OrderRepo.FindOrderById(ctx, id)
	if err != nil {
		return nil, err
	}
	return o, nil
}


func (os *OrderService) FindAllOrders(ctx context.Context) ([]model.Order, error) {
	o, err := os.OrderRepo.FindAllOrders(ctx)
	if err != nil {
		return nil, err
	}
	return o, nil
}


func (os *OrderService) FindOrdersByTenantId(ctx context.Context, tenantId int64) ([]model.Order, error) {
	o, err := os.OrderRepo.FindOrdersByTenantId(ctx, tenantId)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func checkRestaurantAndSubPackageExist(ctx context.Context, rId int64, sId int64, rr *repository.RestaurantRepository, spr *repository.SubPackageRepository) (*model.Restaurant, *model.SubPackage, error) {
	var err error
	var r *model.Restaurant
	r, err = rr.FindRestaurantById(ctx, rId)
	if err != nil {
		return nil, nil, err
	}

	var s *model.SubPackage
	s, err = spr.FindSubPackageById(ctx, sId)
	if s != nil {
		return nil, nil, err
	}

	return r, s, nil
}

func isRestaurantSubPackageMatchTheRequestedPaySubPackage(r *model.Restaurant, sId int64) bool {
	return r.SubPackageId == sId
}

func checkIfRequestedRestaurantBelongToTenant(r *model.Restaurant, tenantId int64) error {
	if r.TenantId != tenantId {
		return appErr.BadRequestError("The requested restaurant does not belong to current tenant.")
	}
	return nil
}
