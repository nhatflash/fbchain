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

type ISubPackageService interface {
	HandleCreateSubPackage(createSubPackageReq *client.CreateSubPackageRequest) (*client.SubPackageResponse, error)
}

type SubPackageService struct {
	Db 			*sql.DB
}

func NewSubPackageService(db *sql.DB) ISubPackageService {
	return &SubPackageService{
		Db: db,
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

	if repository.CheckSubPackageNameExists(name, ss.Db) {
		return nil, appErr.BadRequestError("Subscription package name is already in use.")
	}
	price, err = decimal.NewFromString(priceStr)
	if err != nil {
		return nil, err
	}

	var s *model.SubPackage
	s, err = repository.CreateSubPackage(name, description, durationMonth, price, image, ss.Db)
	if err != nil {
		return nil, err
	}
	return helper.MapToSubPackageResponse(s), nil
}
