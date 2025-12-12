package client

import (
	"time"

	"github.com/nhatflash/fbchain/enum"
	"github.com/shopspring/decimal"
)

type UserResponse struct {
	Id				int64				`json:"id"`
	Email			string				`json:"email"`
	Role			*enum.Role			`json:"role"`
	Phone			string				`json:"phone"`
	Identity		string				`json:"identity"`
	FirstName		string				`json:"firstName"`
	LastName		string				`json:"lastName"`
	Gender			*enum.Gender		`json:"gender"`
	Birthdate		time.Time			`json:"birthdate"`
	PostalCode		string				`json:"postalCode"`
	Address 		string				`json:"address"`
	ProfileImage	string				`json:"profileImage"`
	Status			*enum.UserStatus	`json:"status"`
}

type TenantResponse struct {
	UserId 			int64				`json:"userId"`
	Email			string				`json:"email"`
	Phone			string				`json:"phone"`
	Identity		string				`json:"identity"`
	FirstName		string				`json:"firstName"`
	LastName		string				`json:"lastName"`
	Gender			*enum.Gender		`json:"gender"`
	Birthdate		time.Time			`json:"birthdate"`
	PostalCode		string				`json:"postalCode"`
	Address 		string				`json:"address"`
	ProfileImage	string				`json:"profileImage"`
	Code			string				`json:"code"`
	Description		string				`json:"description"`
	Type 			*enum.TenantType	`json:"type"`
	Notes			string				`json:"notes"`
	Status			*enum.UserStatus	`json:"status"`
}

type SignInResponse struct {
	AccessToken 	string				`json:"accessToken"`
	RefreshToken	string				`json:"refreshToken"`
	LastLogin		time.Time			`json:"lastLogin"`
}


type SubscriptionResponse struct {
	Id				int64				`json:"id"`
	Name			string				`json:"name"`
	Description		string				`json:"description"`
	DurationMonth	int					`json:"durationMonth"`
	Price			decimal.Decimal		`json:"price"`
	IsActive		bool				`json:"isActive"`
	Image			string				`json:"image"`
}