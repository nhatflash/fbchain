package model

import (
	"time"

	"github.com/nhatflash/fbchain/enum"
	"github.com/shopspring/decimal"
)

type User struct {
	Id           int64            `json:"id"`
	Email        string           `json:"email"`
	Password     string           `json:"password"`
	Role         enum.Role        `json:"role"`
	Phone        *string		  `json:"phone"`
	Identity     *string		  `json:"identity"`
	FirstName    string           `json:"firstName"`
	LastName     string           `json:"lastName"`
	Gender       enum.Gender      `json:"gender"`
	Birthdate    time.Time        `json:"birthdate"`
	PostalCode   *string		  `json:"postalCode"`
	Address      *string   		  `json:"address"`
	ProfileImage *string   		  `json:"profileImage"`
	Status       enum.UserStatus  `json:"status"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}

type Tenant struct {
	Id          int64            `json:"id"`
	UserId      int64            `json:"userId"`
	Code        string           `json:"code"`
	Description *string		     `json:"description"`
	Type        enum.TenantType `json:"type"`
	Notes       *string			 `json:"notes"`
}

type Subscription struct {
	Id            int64           `json:"id"`
	Name          string          `json:"name"`
	Description   *string  		  `json:"description"`
	DurationMonth int             `json:"durationMonth"`
	Price         decimal.Decimal `json:"price"`
	IsActive      bool            `json:"isActive"`
	Image         *string  		  `json:"image"`
	Restaurants   []Restaurant
}

type Restaurant struct {
	Id             int64                `json:"id"`
	TenantId       int64                `json:"tenantId"`
	Name           string               `json:"name"`
	Location       string               `json:"location"`
	Description    *string		        `json:"description"`
	ContactEmail   *string		        `json:"contactEmail"`
	ContactPhone   *string		        `json:"contactPhone"`
	PostalCode     string               `json:"postalCode"`
	Type           enum.RestaurantType  `json:"type"`
	AvgRating      decimal.Decimal      `json:"avgRating"`
	IsActive       bool                 `json:"isActive"`
	Notes          *string              `json:"notes"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
	SubscriptionId int64                `json:"subsciptionId"`
	Tenant         *Tenant
	Subscription   *Subscription
	Images         []RestaurantImage
}

type RestaurantImage struct {
	Id           int64     `json:"id"`
	RestaurantId int64     `json:"restaurantId"`
	Image        string    `json:"image"`
	CreatedAt    time.Time `json:"createdAt"`
	Restaurant   *Restaurant
}

type Order struct {
	Id             int64             `json:"id"`
	TenantId       int64             `json:"tenantId"`
	RestaurantId   int64             `json:"restaurantId"`
	SubscriptionId int64             `json:"subscriptionId"`
	OrderDate      time.Time         `json:"orderDate"`
	Status         enum.OrderStatus `json:"status"`
	Amount         decimal.Decimal   `json:"amount"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	Tenant         *Tenant
	Restaurant     *Restaurant
	Subscription   *Subscription
}
