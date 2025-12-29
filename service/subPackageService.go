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
	FindSubPackageById(ctx context.Context, id int64) (*model.SubPackage, error)
	FindAllSubPackages(ctx context.Context) ([]model.SubPackage, error)
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
	var err error
	var exist bool
	exist, err = ss.SubPackageRepo.CheckSubPackageNameExists(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, appErr.BadRequestError("Subscription package name is already in use.")
	}

	var price decimal.Decimal
	price, err = decimal.NewFromString(req.Price)
	if err != nil {
		return nil, err
	}

	var s *model.SubPackage
	s, err = ss.SubPackageRepo.CreateNewSubPackage(ctx, req.Name, req.Description, *req.DurationMonth, price, req.Image)
	if err != nil {
		return nil, err
	}
	return helper.MapToSubPackageResponse(s), nil
}


func (ss *SubPackageService) FindSubPackageById(ctx context.Context, id int64) (*model.SubPackage, error) {
	s, err := ss.SubPackageRepo.FindSubPackageById(ctx, id)
	if err != nil {
		return nil, err
	}
	return s, nil
}


func (ss *SubPackageService) FindAllSubPackages(ctx context.Context) ([]model.SubPackage, error) {
	s, err := ss.SubPackageRepo.FindAllSubPackages(ctx)
	if err != nil {
		return nil, err
	}
	return s, nil
}
