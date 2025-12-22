package service

import (
	"context"

	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
)

type IRestaurantService interface {
	HandleCreateRestaurant(ctx context.Context, req *client.CreateRestaurantRequest, tenantId int64) (*client.RestaurantResponse, error)
	GetRestaurantsByTenantId(tenantId int64) ([]model.Restaurant, error)
	GetAllRestaurants() ([]model.Restaurant, error)
	GetRestaurantById(id int64) (*model.Restaurant, error)
	GetRestaurantImageById(id int64) (*model.RestaurantImage, error)
	GetRestaurantImages(rId int64) ([]model.RestaurantImage, error)
	GetAllRestaurantImages() ([]model.RestaurantImage, error)
}

type RestaurantService struct {
	RestaurantRepo 			*repository.RestaurantRepository
	SubPackageRepo 		*repository.SubPackageRepository
}

func NewRestaurantService(rr *repository.RestaurantRepository, spr *repository.SubPackageRepository) IRestaurantService {
	return &RestaurantService{
		RestaurantRepo: rr,
		SubPackageRepo: spr,
	}
}

func (rs *RestaurantService) HandleCreateRestaurant(ctx context.Context, req *client.CreateRestaurantRequest, tenantId int64) (*client.RestaurantResponse, error) {
	name := req.Name
	location := req.Location
	description := req.Description
	contactEmail := req.ContactEmail
	contactPhone := req.ContactPhone
	postalCode := req.PostalCode
	rType := req.Type
	notes := req.Notes
	images := req.Images

	var err error
	var exist bool
	exist, err = rs.SubPackageRepo.AnySubPackageExists()
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, appErr.NotFoundError("There is no subscription available in the system.")
	}
	if rs.RestaurantRepo.IsRestaurantNameExist(name) {
		return nil, appErr.BadRequestError("Restaurant with this requested name is already exist.")
	}

	var r *model.Restaurant
	r, err = rs.RestaurantRepo.CreateRestaurant(ctx, name, location, description, contactEmail, contactPhone, postalCode, rType, notes, images, tenantId)
	if err != nil {
		return nil, err
	}
	var rImgs []model.RestaurantImage
	rImgs, err = rs.RestaurantRepo.GetRestaurantImages(r.Id)
	if err != nil {
		return nil, err
	}
	return helper.MapToRestaurantResponse(r, rImgs), nil
}

func (rs *RestaurantService) GetRestaurantsByTenantId(tenantId int64) ([]model.Restaurant, error) {
	
	r, err := rs.RestaurantRepo.GetRestaurantsByTenantId(tenantId)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rs *RestaurantService) GetAllRestaurants() ([]model.Restaurant, error) {
	r, err := rs.RestaurantRepo.ListAllRestaurants()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rs *RestaurantService) GetRestaurantById(id int64) (*model.Restaurant, error) {
	r, err := rs.RestaurantRepo.GetRestaurantById(id)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rs *RestaurantService) GetRestaurantImageById(id int64) (*model.RestaurantImage, error) {
	img, err := rs.RestaurantRepo.GetRestaurantImageById(id)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (rs *RestaurantService) GetRestaurantImages(rId int64) ([]model.RestaurantImage, error) {
	imgs, err := rs.RestaurantRepo.GetRestaurantImages(rId)
	if err != nil {
		return nil, err
	}
	return imgs, nil
}


func (rs *RestaurantService) GetAllRestaurantImages() ([]model.RestaurantImage, error) {
	imgs, err := rs.RestaurantRepo.ListAllRestaurantImages()
	if err != nil {
		return nil, err
	}
	return imgs, nil
}