package controller

import (
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/enum"
)

type UserController struct {
	UserRepo 		*repository.UserRepository
}


func NewUserController(ur *repository.UserRepository) *UserController {
	return &UserController{
		UserRepo: ur,
	}
}


func (uc *UserController) ChangeProfile(firstName *string, lastName *string, birthdate *string, gender *enum.Gender, phone *string, identity *string, address *string, postalCode *string, profileImage *string) {
	
}
