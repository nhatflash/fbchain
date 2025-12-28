package client

import (
	"time"

	"github.com/nhatflash/fbchain/enum"
	"github.com/shopspring/decimal"
)

type UserResponse struct {
	Id           int64            `json:"id"`
	Email        string           `json:"email"`
	Role         enum.Role       `json:"role"`
	Phone        *string           `json:"phone"`
	Identity     *string           `json:"identity"`
	FirstName    string           `json:"firstName"`
	LastName     string           `json:"lastName"`
	Gender       enum.Gender     `json:"gender"`
	Birthdate    time.Time        `json:"birthdate"`
	PostalCode   *string           `json:"postalCode"`
	Address      *string           `json:"address"`
	ProfileImage *string           `json:"profileImage"`
	Status       enum.UserStatus `json:"status"`
}

type TenantResponse struct {
	UserId       int64            `json:"userId"`
	Email        string           `json:"email"`
	Phone        *string           `json:"phone"`
	Identity     *string           `json:"identity"`
	FirstName    string           `json:"firstName"`
	LastName     string           `json:"lastName"`
	Gender       enum.Gender     `json:"gender"`
	Birthdate    time.Time        `json:"birthdate"`
	PostalCode   *string           `json:"postalCode"`
	Address      *string           `json:"address"`
	ProfileImage *string           `json:"profileImage"`
	Code         string           `json:"code"`
	Description  *string           `json:"description"`
	Type         enum.TenantType `json:"type"`
	Notes        *string           `json:"notes"`
	Status       enum.UserStatus `json:"status"`
}

type SignInResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	LastLogin    time.Time `json:"lastLogin"`
}

type SubPackageResponse struct {
	Id            int64           `json:"id"`
	Name          string          `json:"name"`
	Description   *string         `json:"description"`
	DurationMonth int             `json:"durationMonth"`
	Price         decimal.Decimal `json:"price"`
	IsActive      bool            `json:"isActive"`
	Image         *string         `json:"image"`
}

type RestaurantResponse struct {
	Id             int64                	`json:"id"`
	TenantId       int64                	`json:"tenantId"`
	Name           string               	`json:"name"`
	Location       string               	`json:"location"`
	Description    *string               	`json:"description"`
	ContactEmail   *string               	`json:"contactEmail"`
	ContactPhone   *string               	`json:"contactPhone"`
	PostalCode     string               	`json:"postalCode"`
	Type           enum.RestaurantType 		`json:"type"`
	AvgRating      decimal.Decimal      	`json:"avgRating"`
	Notes          *string               	`json:"notes"`
	SubPackageId   int64                	`json:"subPackageId"`
	Images         []string             	`json:"images"`
}


type RestaurantImageResponse struct {
	Id 				int64					`json:"id"`
	Image 			string					`json:"image"`
	RestaurantId 	int64					`json:"restaurantId"`
}

type OrderResponse struct {
	Id           int64             	`json:"id"`
	TenantId     int64             	`json:"tenantId"`
	RestaurantId int64             	`json:"restaurantId"`
	OrderDate    time.Time         	`json:"orderDate"`
	Status       enum.OrderStatus 	`json:"status"`
	Amount       decimal.Decimal   	`json:"amount"`
}


type RestaurantItemResponse struct {
	Id 				string				`json:"id"`
	Name 			string 				`json:"name"`
	Description 	*string				`json:"description"`
	Price 			decimal.Decimal 	`json:"price"`
	Type 			enum.ItemType 		`json:"type"`
	Status 			enum.ItemStatus 	`json:"status"`
	Image 			*string				`json:"image"`
	Notes 			*string 			`json:"notes"`
	RestaurantId 	int64				`json:"restaurantId"`
}


type RestaurantTableResponse struct {
	Id 				int64				`json:"id"`
	RestaurantId 	int64				`json:"restaurantId"`
	Label 			string				`json:"label"`
	IsActive 		bool 				`json:"isActive"`
	Notes 			*string				`json:"notes"`
}