package service

import (
	"database/sql"
	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
)

type IRestaurantService interface {
	HandleCreateRestaurant(createRestaurantReq *client.CreateRestaurantRequest, tenantId int64) (*client.RestaurantResponse, error)
	GetRestaurantsByTenantId(tenantId int64) ([]model.Restaurant, error)
	GetAllRestaurants() ([]model.Restaurant, error)
	GetRestaurantById(id int64) (*model.Restaurant, error)
	GetRestaurantImageById(id int64) (*model.RestaurantImage, error)
	GetRestaurantImages(rId int64) ([]model.RestaurantImage, error)
	GetAllRestaurantImages() ([]model.RestaurantImage, error)
}

type RestaurantService struct {
	Db *sql.DB
}

func NewRestaurantService(db *sql.DB) IRestaurantService {
	return &RestaurantService{
		Db: db,
	}
}

func (rs *RestaurantService) HandleCreateRestaurant(createRestaurantReq *client.CreateRestaurantRequest, tenantId int64) (*client.RestaurantResponse, error) {
	name := createRestaurantReq.Name
	location := createRestaurantReq.Location
	description := createRestaurantReq.Description
	contactEmail := createRestaurantReq.ContactEmail
	contactPhone := createRestaurantReq.ContactPhone
	postalCode := createRestaurantReq.PostalCode
	rType := createRestaurantReq.Type
	notes := createRestaurantReq.Notes
	images := createRestaurantReq.Images

	var err error
	var exist bool
	exist, err = repository.AnySubPackageExists(rs.Db)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, appErr.NotFoundError("There is no subscription available in the system.")
	}
	if repository.IsRestaurantNameExist(name, rs.Db) {
		return nil, appErr.BadRequestError("Restaurant with this requested name is already exist.")
	}

	var r *model.Restaurant
	r, err = repository.CreateRestaurant(name, location, description, contactEmail, contactPhone, postalCode, rType, notes, images, tenantId, rs.Db)
	if err != nil {
		return nil, err
	}
	var rImgs []model.RestaurantImage
	rImgs, err = repository.GetRestaurantImages(r.Id, rs.Db)
	if err != nil {
		return nil, err
	}
	return helper.MapToRestaurantResponse(r, rImgs), nil
}

func (rs *RestaurantService) GetRestaurantsByTenantId(tenantId int64) ([]model.Restaurant, error) {
	r, err := repository.GetRestaurantsByTenantId(tenantId, rs.Db)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rs *RestaurantService) GetAllRestaurants() ([]model.Restaurant, error) {
	r, err := repository.ListAllRestaurants(rs.Db)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rs *RestaurantService) GetRestaurantById(id int64) (*model.Restaurant, error) {
	r, err := repository.GetRestaurantById(id, rs.Db)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rs *RestaurantService) GetRestaurantImageById(id int64) (*model.RestaurantImage, error) {
	img, err := repository.GetRestaurantImageById(id, rs.Db)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (rs *RestaurantService) GetRestaurantImages(rId int64) ([]model.RestaurantImage, error) {
	imgs, err := repository.GetRestaurantImages(rId, rs.Db)
	if err != nil {
		return nil, err
	}
	return imgs, nil
}


func (rs *RestaurantService) GetAllRestaurantImages() ([]model.RestaurantImage, error) {
	imgs, err := repository.ListAllRestaurantImages(rs.Db)
	if err != nil {
		return nil, err
	}
	return imgs, nil
}