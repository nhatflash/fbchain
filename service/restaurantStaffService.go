package service

import (
	"github.com/nhatflash/fbchain/repository"
)


type IRestaurantStaffService interface {

}


type RestaurantStaffService struct {
	RestaurantRepo 				*repository.RestaurantRepository
}


func NewRestaurantStaffService() IRestaurantStaffService {
	return &RestaurantStaffService{}
}


