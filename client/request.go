package client

import (
	"github.com/nhatflash/fbchain/enum"
)

type SignInRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type InitialTenantRegisterRequest struct {
	FirstName       string       `json:"firstName" binding:"required,name"`
	LastName        string       `json:"lastName" binding:"required,name"`
	Email           string       `json:"email" binding:"required,email"`
	Password        string       `json:"password" binding:"required"`
	ConfirmPassword string       `json:"confirmPassword" binding:"required"`
	Birthdate       string       `json:"birthdate" binding:"required"`
	Gender          *enum.Gender `json:"gender" binding:"required"`
}

type CompletedTenantRegisterRequest struct {
	Phone        string           `json:"phone" binding:"required,phone"`
	Identity     string           `json:"identity" binding:"required,identity"`
	Address      string           `json:"address" binding:"required"`
	PostalCode   string           `json:"postalCode" binding:"required,postalcode"`
	ProfileImage string           `json:"profileImage"`
	Description  string           `json:"description"`
	Type         *enum.TenantType `json:"type"`
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
	ProfileImage    string           `json:"profileImage"`
	Description     string           `json:"description"`
	Type            *enum.TenantType `json:"type" binding:"required"`
}

type CreateSubPackageRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	DurationMonth int    `json:"durationMonth" binding:"number"`
	Price         string `json:"price" binding:"required,price"`
	Image         string `json:"image"`
}

type CreateRestaurantRequest struct {
	Name         string               `json:"name" binding:"required"`
	Location     string               `json:"location" binding:"required"`
	Description  string               `json:"description"`
	ContactEmail string               `json:"contactEmail" binding:"email"`
	ContactPhone string               `json:"contactPhone" binding:"phone"`
	PostalCode   string               `json:"postalCode" binding:"required,postalcode"`
	Type         *enum.RestaurantType `json:"type" binding:"required"`
	Notes        string               `json:"notes" binding:"required"`
	Images       []string             `json:"image"`
}

type PaySubscriptionRequest struct {
	RestaurantId   int64 `json:"restaurantId" binding:"required"`
	SubscriptionId int64 `json:"subscriptionId" binding:"required"`
}
