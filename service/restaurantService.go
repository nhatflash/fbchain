package service

import (
	"context"

	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/shopspring/decimal"
)

type IRestaurantService interface {
	HandleCreateRestaurant(ctx context.Context, req *client.CreateRestaurantRequest, tenantId int64) (*client.RestaurantResponse, error)
	GetRestaurantsByTenantId(ctx context.Context, tenantId int64) ([]model.Restaurant, error)
	GetAllRestaurants(ctx context.Context) ([]model.Restaurant, error)
	GetRestaurantById(ctx context.Context, id int64) (*model.Restaurant, error)
	GetRestaurantImageById(ctx context.Context, id int64) (*model.RestaurantImage, error)
	GetRestaurantImages(ctx context.Context, restaurantId int64) ([]model.RestaurantImage, error)
	GetAllRestaurantImages(ctx context.Context) ([]model.RestaurantImage, error)
	HandleAddNewRestaurantItem(ctx context.Context, restaurantId int64, tenantId int64, req *client.AddRestaurantItemRequest) (*client.RestaurantItemResponse, error)
	GetItemsByRestaurantId(ctx context.Context, restaurantId int64) ([]model.RestaurantItem, error)
	GetAllRestaurantItems(ctx context.Context) ([]model.RestaurantItem, error)
	GetRestaurantItemById(ctx context.Context, id int64) (*model.RestaurantItem, error)
}

type RestaurantService struct {
	RestaurantRepo 		*repository.RestaurantRepository
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
	
	if err = validateCreateRestaurantRequest(ctx, name, rs.SubPackageRepo, rs.RestaurantRepo); err != nil {
		return nil, err
	}
	var s *model.SubPackage
	s, err = rs.SubPackageRepo.GetFirstSubPackage(ctx)
	if err != nil {
		return nil, err
	}
	var r *model.Restaurant
	r, err = rs.RestaurantRepo.CreateRestaurant(ctx, name, location, description, contactEmail, contactPhone, postalCode, *rType, notes, s.Id, images, tenantId)
	if err != nil {
		return nil, err
	}
	var rImgs []model.RestaurantImage
	rImgs, err = rs.RestaurantRepo.GetRestaurantImages(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return helper.MapToRestaurantResponse(r, rImgs), nil
}



func (rs *RestaurantService) GetRestaurantsByTenantId(ctx context.Context, tenantId int64) ([]model.Restaurant, error) {
	
	r, err := rs.RestaurantRepo.GetRestaurantsByTenantId(ctx, tenantId)
	if err != nil {
		return nil, err
	}
	return r, nil
}



func (rs *RestaurantService) GetAllRestaurants(ctx context.Context) ([]model.Restaurant, error) {
	r, err := rs.RestaurantRepo.ListAllRestaurants(ctx)
	if err != nil {
		return nil, err
	}
	return r, nil
}




func (rs *RestaurantService) GetRestaurantById(ctx context.Context, id int64) (*model.Restaurant, error) {
	r, err := rs.RestaurantRepo.GetRestaurantById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r, nil
}




func (rs *RestaurantService) GetRestaurantImageById(ctx context.Context, id int64) (*model.RestaurantImage, error) {
	img, err := rs.RestaurantRepo.GetRestaurantImageById(ctx, id)
	if err != nil {
		return nil, err
	}
	return img, nil
}



func (rs *RestaurantService) GetRestaurantImages(ctx context.Context, restaurantId int64) ([]model.RestaurantImage, error) {
	imgs, err := rs.RestaurantRepo.GetRestaurantImages(ctx, restaurantId)
	if err != nil {
		return nil, err
	}
	return imgs, nil
}


func (rs *RestaurantService) GetAllRestaurantImages(ctx context.Context) ([]model.RestaurantImage, error) {
	imgs, err := rs.RestaurantRepo.ListAllRestaurantImages(ctx)
	if err != nil {
		return nil, err
	}
	return imgs, nil
}



func (rs *RestaurantService) HandleAddNewRestaurantItem(ctx context.Context, restaurantId int64, tenantId int64, req *client.AddRestaurantItemRequest) (*client.RestaurantItemResponse, error) {
	var err error 
	var r *model.Restaurant
	
	r, err = rs.RestaurantRepo.GetRestaurantById(ctx, restaurantId)
	if err != nil {
		return nil, err
	}

	if r.TenantId != tenantId {
		return nil, appErr.UnauthorizedError("You are not allowed to add new item on this restaurant.")
	}

	var price decimal.Decimal
	price, err = decimal.NewFromString(req.Price)
	if err != nil {
		return nil, appErr.BadRequestError("Invalid price.")
	}

	var i *model.RestaurantItem
	i, err = rs.RestaurantRepo.AddNewRestaurantItem(ctx, req.Name, req.Description, price, req.Type, req.Image, req.Notes, restaurantId)
	if err != nil {
		return nil, err
	}
	return helper.MapToRestaurantItemResponse(i), nil
}



func (rs *RestaurantService) GetItemsByRestaurantId(ctx context.Context, restaurantId int64) ([]model.RestaurantItem, error) {
	items, err := rs.RestaurantRepo.GetItemsByRestaurantId(ctx, restaurantId)
	if err != nil {
		return nil, err
	}
	return items, nil
}


func (rs *RestaurantService) GetAllRestaurantItems(ctx context.Context) ([]model.RestaurantItem, error) {
	items, err := rs.RestaurantRepo.GetAllRestaurantItems(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}


func (rs *RestaurantService) GetRestaurantItemById(ctx context.Context, id int64) (*model.RestaurantItem, error) {
	item, err := rs.RestaurantRepo.GetRestaurantItemById(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}




func validateCreateRestaurantRequest(ctx context.Context, name string, subPackageRepo *repository.SubPackageRepository, resRepo *repository.RestaurantRepository) error {
	var err error
	var exist bool
	exist, err = subPackageRepo.AnySubPackageExists(ctx)
	if err != nil {
		return err
	}
	if !exist {
		return appErr.NotFoundError("There is no subscription available in the system.")
	}
	exist, err = resRepo.IsRestaurantNameExist(ctx, name)
	if err != nil {
		return err
	}
	if exist {
		return appErr.BadRequestError("Restaurant with this requested name is already exist.")
	}
	return nil
}