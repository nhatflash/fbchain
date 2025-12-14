package service

import (
	"database/sql"

	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/shopspring/decimal"
)

func HandleCreateSubscription(createSubScriptionReq *client.CreateSubscriptionRequest, db *sql.DB) (*client.SubscriptionResponse, error) {
	name := createSubScriptionReq.Name
	description := createSubScriptionReq.Description
	durationMonth := createSubScriptionReq.DurationMonth
	priceStr := createSubScriptionReq.Price
	image := createSubScriptionReq.Image

	var err error
	var price decimal.Decimal

	if repository.CheckSubscriptionNameExists(name, db) {
		return nil, appErr.BadRequestError("Subscription name is already in use.")
	}
	price, err = decimal.NewFromString(priceStr)
	if err != nil {
		return nil, err
	}

	var s *model.Subscription
	s, err = repository.CreateSubscription(name, description, durationMonth, price, image, db)
	if err != nil {
		return nil, err
	}
	return helper.MapToSubscriptionResponse(s), nil
}
