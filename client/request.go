package client

import (
	"github.com/nhatflash/fbchain/enum"
)

type LoginRequest struct {
	Login			string				`json:"login" binding:"required"`
	Password		string				`json:"password" binding:"required"`
}

type RegisterTenantRequest struct {
	FirstName		string				`json:"firstName" binding:"required"`
	LastName		string				`json:"lastName" binding:"required"`
	Email			string 				`json:"email" binding:"required,email"`
	Password 		string				`json:"password" binding:"required"`
	ConfirmPassword	string				`json:"confirmPassword" binding:"required"`
	Birthdate		string				`json:"birthdate" binding:"required"`
	Gender			enum.Gender			`json:"gender" binding:"required"`	
}