package service

import (
	"database/sql"

	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
)

func HandlePaySubscription(paySubscriptionReq *client.PaySubscriptionRequest, tenantId int64, db *sql.DB) (*client.OrderResponse, error) {
	restaurantId := paySubscriptionReq.RestaurantId
	subscriptionId := paySubscriptionReq.SubscriptionId

	var err error
	var r *model.Restaurant
	var s *model.Subscription
	r, s, err = checkRestaurantAndSubscriptionExist(restaurantId, subscriptionId, db)
	if err != nil {
		return nil, err
	}
	err = checkIfRequestedRestaurantBelongToTenant(r, tenantId)
	if err != nil {
		return nil, err
	}
	if isRestaurantSubscriptionMatchTheRequestedPaySubscription(r, s.Id) {
		return nil, appErr.BadRequestError("The requested subscription is already registered on this restaurant.")
	}
	err = repository.CreateInitialOrder(restaurantId, subscriptionId, &s.Price, tenantId, db)
	if err != nil {
		return nil, err
	}
	var order *model.Order
	order, err = repository.GetLatestTenantOrder(tenantId, db)
	if err != nil {
		return nil, err
	}
	return helper.MapToOrderResponse(order), nil
}

func checkRestaurantAndSubscriptionExist(rId int64, sId int64, db *sql.DB) (*model.Restaurant, *model.Subscription, error) {
	var err error
	var r *model.Restaurant
	r, err = repository.GetRestaurantById(rId, db)
	if err != nil {
		return nil, nil, err
	}

	var s *model.Subscription
	s, err = repository.GetSubscriptionById(sId, db)
	if s != nil {
		return nil, nil, err
	}

	return r, s, nil
}

func isRestaurantSubscriptionMatchTheRequestedPaySubscription(r *model.Restaurant, sId int64) bool {
	return r.SubscriptionId == sId
}

func checkIfRequestedRestaurantBelongToTenant(r *model.Restaurant, tenantId int64) error {
	if r.TenantId != tenantId {
		return appErr.BadRequestError("The requested restaurant does not belong to current tenant.")
	}
	return nil
}
