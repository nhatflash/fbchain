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

type ISubPackageService interface {
	HandleCreateSubPackage(ctx context.Context, req *client.CreateSubPackageRequest) (*client.SubPackageResponse, error)
}

type SubPackageService struct {
	SubPackageRepo 			*repository.SubPackageRepository
}

func NewSubPackageService(spr *repository.SubPackageRepository) ISubPackageService {
	return &SubPackageService{
		SubPackageRepo: spr,
	}
}

func (ss *SubPackageService) HandleCreateSubPackage(ctx context.Context, req *client.CreateSubPackageRequest) (*client.SubPackageResponse, error) {
	name := req.Name
	description := req.Description
	durationMonth := req.DurationMonth
	priceStr := req.Price
	image := req.Image

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
	s, err = ss.SubPackageRepo.CreateSubPackage(ctx, name, description, *durationMonth, price, image)
	if err != nil {
		return nil, err
	}
	return helper.MapToSubPackageResponse(s), nil
}
