package service

import (
	"github.com/nhatflash/fbchain/client"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/shopspring/decimal"
)

type ISubPackageService interface {
	HandleCreateSubPackage(createSubPackageReq *client.CreateSubPackageRequest) (*client.SubPackageResponse, error)
}

type SubPackageService struct {
	SubPackageRepo 			*repository.SubPackageRepository
}

func NewSubPackageService(spr *repository.SubPackageRepository) ISubPackageService {
	return &SubPackageService{
		SubPackageRepo: spr,
	}
}

func (ss *SubPackageService) HandleCreateSubPackage(createSubPackageReq *client.CreateSubPackageRequest) (*client.SubPackageResponse, error) {
	name := createSubPackageReq.Name
	description := createSubPackageReq.Description
	durationMonth := createSubPackageReq.DurationMonth
	priceStr := createSubPackageReq.Price
	image := createSubPackageReq.Image

	var err error
	var price decimal.Decimal

	if ss.SubPackageRepo.CheckSubPackageNameExists(name) {
		return nil, appErr.BadRequestError("Subscription package name is already in use.")
	}
	price, err = decimal.NewFromString(priceStr)
	if err != nil {
		return nil, err
	}

	var s *model.SubPackage
	s, err = ss.SubPackageRepo.CreateSubPackage(name, description, durationMonth, price, image)
	if err != nil {
		return nil, err
	}
	return helper.MapToSubPackageResponse(s), nil
}
