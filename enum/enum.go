package enum

type Role string

const (
	ADMIN Role = "ADMIN"
	MANAGER Role = "MANAGER"
	STAFF Role = "STAFF"
	TENANT Role = "TENANT"
	RESTAURANT_STAFF = "RESTAURANT_STAFF"
)


type Gender string

const (
	MALE Gender = "MALE"
	FEMALE Gender = "FEMALE"
)


type UserStatus string

const (
	ACTIVE UserStatus = "ACTIVE"
	INACTIVE UserStatus = "INACTIVE"
	LOCKED UserStatus = "LOCKED"
	DELETED UserStatus = "DELETED"
)


type TenantType string

const (
	FREE TenantType = "FREE"
	PERSONAL TenantType = "PERSONAL"
	BUSINESS TenantType = "BUSINESS"
)



type RestaurantType string

const (
	FOOD TenantType = "FOOD"
	BEVERAGE TenantType = "BEVERAGE"
	FOOD_AND_BEVERAGE TenantType = "FB"
)

