package client

import (
	"github.com/nhatflash/fbchain/enum"
)

type LoginRequest struct {
	Login			string				`json:"login" binding:"required"`
	Password		string				`json:"password" binding:"required"`
}

type InitializedTenantRegisterRequest struct {
	FirstName		string				`json:"firstName" binding:"required,name"`
	LastName		string				`json:"lastName" binding:"required,name"`
	Email			string 				`json:"email" binding:"required,email"`
	Password 		string				`json:"password" binding:"required"`
	ConfirmPassword	string				`json:"confirmPassword" binding:"required"`
	Birthdate		string				`json:"birthdate" binding:"required"`
	Gender			enum.Gender			`json:"gender" binding:"required"`	
}

type CompletedTenantRegisterRequest struct {
	Phone 			string				`json:"phone" binding:"required,phone"`
	Identity		string 				`json:"identity" binding:"required,identity"`
	Address			string				`json:"address" binding:"required"`
	PostalCode		string				`json:"postalCode" binding:"required,postalcode"`
	ProfileImage	string				`json:"profileImage"`
}