package model

import (
	"time"

	"github.com/nhatflash/fbchain/enum"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	IsVerified	 bool 			  `json:"isVerified"`
}

type Tenant struct {
	Id          int64            `json:"id"`
	UserId      int64            `json:"userId"`
	Code        string           `json:"code"`
	Description *string		     `json:"description"`
	Type        enum.TenantType `json:"type"`
	Notes       *string			 `json:"notes"`
}

type SubPackage struct {
	Id            int64           `json:"id"`
	Name          string          `json:"name"`
	Description   *string  		  `json:"description"`
	DurationMonth int             `json:"durationMonth"`
	Price         decimal.Decimal `json:"price"`
	IsActive      bool            `json:"isActive"`
	CreatedAt	  time.Time 	  `json:"createdAt"`
	UpdatedAt 	  time.Time 	  `json:"updatedAt"`
	Image         *string  		  `json:"image"`
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
	SubPackageId   int64                  `json:"subPackageId"`
}

type RestaurantImage struct {
	Id           int64     				`json:"id"`
	RestaurantId int64     				`json:"restaurantId"`
	Image        string    				`json:"image"`
	CreatedAt    time.Time 				`json:"createdAt"`
}

type Order struct {
	Id             int64             `json:"id"`
	TenantId       int64             `json:"tenantId"`
	RestaurantId   int64             `json:"restaurantId"`
	SubPackageId   int64             `json:"subPackageId"`
	OrderDate      time.Time         `json:"orderDate"`
	Status         enum.OrderStatus  `json:"status"`
	Amount         decimal.Decimal   `json:"amount"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}


type Payment struct {
	Id 				int64				`json:"id"`
	OrderId			int64				`json:"orderId"`
	Amount			decimal.Decimal 	`json:"amount"`
	Method			enum.PaymentMethod 	`json:"method"`
	BankCode 		*string				`json:"bankCode"`
	Status			enum.PaymentStatus	`json:"status"`
	PaymentDate		time.Time			`json:"paymentDate"`
	Notes 			*string				`json:"notes"`
}


type RestaurantItem struct {
	Id 				bson.ObjectID		`json:"id" bson:"_id,omitempty"`
	Name 			string 				`json:"name" bson:"name"`
	Description 	*string 			`json:"description" bson:"description"`
	Price 			bson.Decimal128		`json:"price" bson:"price"`
	Type			enum.ItemType 		`json:"type" bson:"type"`
	Status 			enum.ItemStatus 	`json:"status" bson:"status"`
	Image 			*string 			`json:"image" bson:"image"`
	Notes 			*string 			`json:"notes" bson:"notes"`
	CreatedAt 		time.Time 			`json:"createdAt" bson:"createdAt"`
	UpdatedAt 		time.Time 			`json:"updatedAt" bson:"updatedAt"`
	RestaurantId 	int64				`json:"restaurantId" bson:"restaurantId"`
}


type RestaurantTable struct {
	Id 				int64				`json:"id"`
	RestaurantId 	int64				`json:"restaurantId"`
	Label 			string 				`json:"label"`
	IsActive 		bool				`json:"isActive"`
	Notes 			*string				`json:"notes"`
	CreatedAt 		time.Time			`json:"createdAt"`
}


type RestaurantOrder struct {
	Id 				int64						`json:"id"`
	RestaurantId 	int64						`json:"restaurantId"`
	TableId 		int64						`json:"tableId"`
	Amount 			decimal.Decimal 			`json:"amount"`
	Status 			enum.RestaurantOrderStatus 	`json:"status"`
	Notes 			*string 					`json:"notes"`
	CreatedAt 		time.Time 					`json:"createdAt"`
	UpdatedAt 		time.Time 					`json:"updatedAt"`
	Items 			[]RestaurantOrderItem 		`json:"items"`
}


type RestaurantOrderItem struct {
	Id 				int64 						`json:"id"`
	ROrderId 		int64						`json:"rOrderId"`
	ItemId 			string 						`json:"itemId"`
	Quantity 		int 						`json:"quantity"`
	Total 			decimal.Decimal 			`json:"total"`
}


type RestaurantPayment struct {
	Id 				int64						`json:"id"`
	ROrderId 		int64 						`json:"rOrderId"`
	Amount 			decimal.Decimal 			`json:"amount"`
	BankCode 		*string						`json:"bankCode"`
	Method 			enum.PaymentMethod   		`json:"method"`
	Status 			enum.PaymentStatus 			`json:"status"`
	IsCashed 		bool 						`json:"isCashed"`
	CreatedAt 		time.Time 					`json:"createdAt"`
}