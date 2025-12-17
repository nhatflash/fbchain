package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.
import (
	"github.com/nhatflash/fbchain/service"
)

type Resolver struct{
	UserService		service.IUserService
	TenantService	service.ITenantService
	RestaurantService	service.IRestaurantService
}
