package service

import (
	"database/sql"

	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/repository"
	"github.com/shopspring/decimal"
)

func HandleCreateSubscription(createSubScriptionReq *client.CreateSubscriptionRequest, db *sql.DB) (*client.SubscriptionResponse, error) {
	name := createSubScriptionReq.Name
	description := createSubScriptionReq.Description
	durationMonth := createSubScriptionReq.DurationMonth
	priceStr := createSubScriptionReq.Price
	image := createSubScriptionReq.Image

	if repository.CheckSubscriptionNameExists(name, db) {
		return nil, appErr.BadRequestError("Subscription name is already in use.")
	}
	price, pErr := decimal.NewFromString(priceStr)
	if pErr != nil {
		return nil, pErr
	}
	subscription, createErr := repository.CreateSubscription(name, description, durationMonth, price, image, db)
	if createErr != nil {
		return nil, createErr
	}
	return helper.MapToSubscriptionResponse(subscription), nil
} 