package service

import (
	"database/sql"

	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/repository"
)

type IRestaurantService interface {
	HandleCreateRestaurant(createRestaurantReq *client.CreateRestaurantRequest, tenantId int64, db *sql.DB) (*client.RestaurantResponse, error)
}

type RestaurantService struct {
	Db 				*sql.DB
}

func NewRestaurantService(db *sql.DB) IRestaurantService {
	return &RestaurantService{
		Db: db,
	}
}


func (*RestaurantService) HandleCreateRestaurant(createRestaurantReq *client.CreateRestaurantRequest, tenantId int64, db *sql.DB) (*client.RestaurantResponse, error) {
	name := createRestaurantReq.Name
	location := createRestaurantReq.Location
	description := createRestaurantReq.Description
	contactEmail := createRestaurantReq.ContactEmail
	contactPhone := createRestaurantReq.ContactPhone
	postalCode := createRestaurantReq.PostalCode
	rType := createRestaurantReq.Type
	notes := createRestaurantReq.Notes
	images := createRestaurantReq.Images

	sExist, sExistErr := repository.AnySubscriptionExists(db)
	if sExistErr != nil {
		return nil, sExistErr
	}
	if !sExist {
		return nil, appErr.NotFoundError("There is no subscription available in the system.")
	}
	if repository.IsRestaurantNameExist(name, db) {
		return nil, appErr.BadRequestError("Restaurant with this requested name is already exist.")
	}

	r, rErr := repository.CreateRestaurant(name, location, description, contactEmail, contactPhone, postalCode, rType, notes, images, tenantId, db)
	if rErr != nil {
		return nil, rErr
	}
	return helper.MapToRestaurantResponse(r), nil
}