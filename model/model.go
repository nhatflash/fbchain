package model

import (
	"database/sql"
	"time"

	"github.com/nhatflash/fbchain/enum"
	"github.com/shopspring/decimal"
)

type User struct {
	Id           int64            `json:"id"`
	Email        string           `json:"email"`
	Password     string           `json:"password"`
	Role         *enum.Role       `json:"role"`
	Phone        sql.NullString   `json:"phone"`
	Identity     sql.NullString   `json:"identity"`
	FirstName    string           `json:"firstName"`
	LastName     string           `json:"lastName"`
	Gender       *enum.Gender     `json:"gender"`
	Birthdate    time.Time        `json:"birthdate"`
	PostalCode   sql.NullString   `json:"postalCode"`
	Address      sql.NullString   `json:"address"`
	ProfileImage sql.NullString   `json:"profileImage"`
	Status       *enum.UserStatus `json:"status"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}

type Tenant struct {
	Id          int64            `json:"id"`
	UserId      int64            `json:"userId"`
	Code        string           `json:"code"`
	Description sql.NullString   `json:"description"`
	Type        *enum.TenantType `json:"type"`
	Notes       sql.NullString   `json:"notes"`
	User        *User
	Restaurants []Restaurant
}

type Subscription struct {
	Id            int64           `json:"id"`
	Name          string          `json:"name"`
	Description   sql.NullString  `json:"description"`
	DurationMonth int             `json:"durationMonth"`
	Price         decimal.Decimal `json:"price"`
	IsActive      bool            `json:"isActive"`
	Image         sql.NullString  `json:"image"`
	Restaurants   []Restaurant
}

type Restaurant struct {
	Id             int64                `json:"id"`
	TenantId       int64                `json:"tenantId"`
	Name           string               `json:"name"`
	Location       string               `json:"location"`
	Description    sql.NullString       `json:"description"`
	ContactEmail   sql.NullString       `json:"contactEmail"`
	ContactPhone   sql.NullString       `json:"contactPhone"`
	PostalCode     string               `json:"postalCode"`
	Type           *enum.RestaurantType `json:"type"`
	AvgRating      decimal.Decimal      `json:"avgRating"`
	IsActive       bool                 `json:"isActive"`
	Notes          sql.NullString       `json:"notes"`
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
	Status         *enum.OrderStatus `json:"status"`
	Amount         decimal.Decimal   `json:"amount"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	Tenant         *Tenant
	Restaurant     *Restaurant
	Subscription   *Subscription
}
