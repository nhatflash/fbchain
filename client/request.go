package client

import (
	"github.com/nhatflash/fbchain/enum"
)

type SignInRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TenantSignUpRequest struct {
	FirstName       string           `json:"firstName" binding:"required,name"`
	LastName        string           `json:"lastName" binding:"required,name"`
	Email           string           `json:"email" binding:"required,email"`
	Password        string           `json:"password" binding:"required"`
	ConfirmPassword string           `json:"confirmPassword" binding:"required"`
	Birthdate       string           `json:"birthdate" binding:"required"`
	Gender          *enum.Gender     `json:"gender" binding:"required"`
	Phone           string           `json:"phone" binding:"required,phone"`
	Identity        string           `json:"identity" binding:"required,identity"`
	Address         string           `json:"address" binding:"required"`
	PostalCode      string           `json:"postalCode" binding:"required,postalcode"`
	ProfileImage    *string           `json:"profileImage" binding:"omitempty"`
	Description     *string           `json:"description" binding:"omitempty"`
	Type            *enum.TenantType `json:"type" binding:"required"`
}

type CreateSubPackageRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   *string `json:"description" binding:"omitempty"`
	DurationMonth *int    `json:"durationMonth" binding:"required,number"`
	Price         string `json:"price" binding:"required,price"`
	Image         *string `json:"image" binding:"omitempty"`
}

type CreateRestaurantRequest struct {
	Name         string               `json:"name" binding:"required"`
	Location     string               `json:"location" binding:"required"`
	Description  *string               `json:"description" binding:"omitempty"`
	ContactEmail *string               `json:"contactEmail" binding:"omitempty,email"`
	ContactPhone *string               `json:"contactPhone" binding:"omitempty,phone"`
	PostalCode   string               `json:"postalCode" binding:"required,postalcode"`
	Type         *enum.RestaurantType `json:"type" binding:"required"`
	Notes        string               `json:"notes" binding:"required"`
	Images       []string             `json:"image"`
}

type PaySubPackageRequest struct {
	RestaurantId   *int64 `json:"restaurantId" binding:"required"`
	SubPackageId   *int64 `json:"subPackageId" binding:"required"`
}


type UpdateProfileRequest struct {
	FirstName 		*string				`json:"firstName" binding:"omitempty,name"`
	LastName 		*string				`json:"lastName" binding:"omitempty,name"`
	Phone			*string				`json:"phone" binding:"omitempty,phone"`
	Identity 		*string				`json:"identity" binding:"omitempty,identity"`
	Gender 			*enum.Gender		`json:"gender" binding:"omitempty"`
	Birthdate		*string				`json:"birthdate" binding:"omitempty"`
	PostalCode		*string				`json:"postalCode" binding:"omitempty,postalcode"`
	Address 		*string				`json:"address" binding:"omitempty"`
	ProfileImage	*string				`json:"profileImage" binding:"omitempty"`
}

type VerifyChangePasswordRequest struct {
	VerifiedCode 	string			`json:"verifiedCode" binding:"required"`
}

type ChangePasswordRequest struct {
	NewPassword			string				`json:"newPassword" binding:"required"`
	ConfirmNewPassword	string				`json:"confirmNewPassword" binding:"required"`
}