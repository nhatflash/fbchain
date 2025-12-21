package enum

type Role string

const (
	ROLE_ADMIN            Role = "ADMIN"
	ROLE_MANAGER          Role = "MANAGER"
	ROLE_STAFF            Role = "STAFF"
	ROLE_TENANT           Role = "TENANT"
	ROLE_RESTAURANT_STAFF      = "RESTAURANT_STAFF"
)

type Gender string

const (
	GENDER_MALE   Gender = "MALE"
	GENDER_FEMALE Gender = "FEMALE"
)

type UserStatus string

const (
	USER_ACTIVE   UserStatus = "ACTIVE"
	USER_INACTIVE UserStatus = "INACTIVE"
	USER_PENDING  UserStatus = "PENDING"
	USER_LOCKED   UserStatus = "LOCKED"
	USER_DELETED  UserStatus = "DELETED"
)

type TenantType string

const (
	TENANT_PERSONAL TenantType = "PERSONAL"
	TENANT_BUSINESS TenantType = "BUSINESS"
)

type RestaurantType string

const (
	RESTAURANT_FOOD              TenantType = "FOOD"
	RESTAURANT_BEVERAGE          TenantType = "BEVERAGE"
	RESTAURANT_FOOD_AND_BEVERAGE TenantType = "FB"
)

type OrderStatus string

const (
	ORDER_PENDING   OrderStatus = "PENDING"
	ORDER_COMPLETED OrderStatus = "COMPLETED"
	ORDER_CANCELED  OrderStatus = "CANCELED"
)


type PaymentStatus string

const (
	PAYMENT_SUCCESS PaymentStatus = "SUCCESS"
	PAYMENT_FAILED PaymentStatus = "FAILED"
)


type PaymentMethod string 

const (
	PAYMENT_CASH PaymentMethod = "CASH"
	PAYMENT_VNPAY PaymentMethod = "VNPAY"
)