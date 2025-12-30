package graph

import (
	"strconv"

	gqlModel "github.com/nhatflash/fbchain/graph/model"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/scalar"
)

func MapToGqlUser(u *model.User) *gqlModel.User {
	return &gqlModel.User{
		ID:           strconv.FormatInt(u.Id, 10),
		Email:        u.Email,
		Phone:        u.Phone,
		Identity:     u.Identity,
		FirstName:    &u.FirstName,
		LastName:     &u.LastName,
		Gender:       &u.Gender,
		Birthdate:    (*scalar.CustomDate)(&u.Birthdate),
		PostalCode:   u.PostalCode,
		Address:      u.Address,
		ProfileImage: u.ProfileImage,
	}
}


func MapToGqlUsers(users []model.User) []*gqlModel.User {
	var gqlUsers []*gqlModel.User
	for _, u := range users {
		gqlUser := MapToGqlUser(&u)
		gqlUsers = append(gqlUsers, gqlUser)
	}
	return gqlUsers
}



func MapToGqlTenant(t *model.Tenant) *gqlModel.Tenant {
	idStr := strconv.FormatInt(t.UserId, 10)
	return &gqlModel.Tenant{
		ID: strconv.FormatInt(t.Id, 10),
		Code: t.Code,
		Description: t.Description,
		Type: &t.Type,
		Notes: t.Notes,
		UserID: &idStr,
	}
}


func MapToGqlTenants(tenants []model.Tenant) []*gqlModel.Tenant {
	var gqlTenants []*gqlModel.Tenant
	for _, t := range tenants {
		gqlTenant := MapToGqlTenant(&t)
		gqlTenants = append(gqlTenants, gqlTenant)
	}
	return gqlTenants
}


func MapToGqlRestaurant(r *model.Restaurant) *gqlModel.Restaurant {
	idStr := strconv.FormatInt(r.TenantId, 10)
	avgRating := r.AvgRating.String()
	return &gqlModel.Restaurant{
		ID: strconv.FormatInt(r.Id, 10),
		Name: r.Name,
		Location: &r.Location,
		Description: r.Description,
		ContactEmail: r.ContactEmail,
		ContactPhone: r.ContactPhone,
		PostalCode: &r.PostalCode,
		Type: &r.Type,
		AvgRating: &avgRating,
		IsActive: &r.IsActive,
		Notes: r.Notes,
		TenantID: &idStr,
	}
}


func MapToGqlRestaurants(restaurants []model.Restaurant) []*gqlModel.Restaurant {
	var gqlRestaurants []*gqlModel.Restaurant
	for _, r := range restaurants {
		gqlRestaurant := MapToGqlRestaurant(&r)
		gqlRestaurants = append(gqlRestaurants, gqlRestaurant)
	}
	return gqlRestaurants
}


func MapToGqlRestaurantImage(rImg *model.RestaurantImage) *gqlModel.RestaurantImage {
	idStr := strconv.FormatInt(rImg.RestaurantId, 10)
	return &gqlModel.RestaurantImage{
		ID: strconv.FormatInt(rImg.Id, 10),
		Image: rImg.Image,
		CreatedAt: &rImg.CreatedAt,
		RestaurantID: &idStr,
	}
}


func MapToGqlRestaurantImages(rImgs []model.RestaurantImage) []*gqlModel.RestaurantImage {
	var gqlRImages []*gqlModel.RestaurantImage
	for _, i := range rImgs {
		gqlRImage := MapToGqlRestaurantImage(&i)
		gqlRImages = append(gqlRImages, gqlRImage)
	}
	return gqlRImages
}


func MapToGqlRestaurantItem(item *model.RestaurantItem) *gqlModel.RestaurantItem {
	idStr := strconv.FormatInt(item.RestaurantId, 10)
	price := item.Price.String()
	itemId := item.Id.Hex()
	return &gqlModel.RestaurantItem{
		ID: itemId,
		Name: item.Name,
		Description: item.Description,
		Price: &price,
		Type: &item.Type,
		Status: &item.Status,
		Image: item.Image,
		Notes: item.Notes,
		RestaurantID: &idStr,
	}
}


func MapToGqlRestaurantItems(items []model.RestaurantItem) []*gqlModel.RestaurantItem {
	var gqlRItems []*gqlModel.RestaurantItem
	for _, i := range items {
		gqlRItem := MapToGqlRestaurantItem(&i)
		gqlRItems = append(gqlRItems, gqlRItem)
	}
	return gqlRItems
}


func MapToGqlRestaurantTable(table *model.RestaurantTable) *gqlModel.RestaurantTable {
	tableId := strconv.FormatInt(table.Id, 10)
	restaurantId := strconv.FormatInt(table.RestaurantId, 10)
	return &gqlModel.RestaurantTable{
		ID: tableId,
		Label: table.Label,
		RestaurantID: &restaurantId,
		IsActive: &table.IsActive,
		Notes: table.Notes,
		CreatedAt: &table.CreatedAt,
	}
}


func MapToGqlRestaurantTables(tables []model.RestaurantTable) []*gqlModel.RestaurantTable {
	var gqlTables []*gqlModel.RestaurantTable
	for _, t := range tables {
		gqlTable := MapToGqlRestaurantTable(&t)
		gqlTables = append(gqlTables, gqlTable)
	}
	return gqlTables
}


func MapToGqlOrder(o *model.Order) *gqlModel.Order {
	orderId := strconv.FormatInt(o.Id, 10)
	tenantId := strconv.FormatInt(o.TenantId, 10)
	restaurantId := strconv.FormatInt(o.RestaurantId, 10)
	subPackageId := strconv.FormatInt(o.SubPackageId, 10)
	return &gqlModel.Order{
		ID: orderId,
		Amount: o.Amount.String(),
		OrderDate: &o.OrderDate,
		Status: &o.Status,
		TenantID: &tenantId,
		RestaurantID: &restaurantId,
		SubPackageID: &subPackageId,
		UpdatedAt: &o.UpdatedAt,
	}
}


func MapToGqlOrders(orders []model.Order) []*gqlModel.Order {
	var gqlOrders []*gqlModel.Order
	for _, o := range orders {
		gqlOrder := MapToGqlOrder(&o)
		gqlOrders = append(gqlOrders, gqlOrder)
	}
	return gqlOrders
}


func MapToGqlSubPackage(s *model.SubPackage) *gqlModel.SubPackage {
	subId := strconv.FormatInt(s.Id, 10)
	durationMonth := int32(s.DurationMonth)
	price := s.Price.String()
	return &gqlModel.SubPackage{
		ID: subId,
		Name: s.Name,
		Description: s.Description,
		DurationMonth: &durationMonth,
		Price: &price,
		IsActive: &s.IsActive,
		CreatedAt: &s.CreatedAt,
		UpdatedAt: &s.UpdatedAt,
		Image: s.Image,
	}
}


func MapToGqlSubPackages(subPackages []model.SubPackage) []*gqlModel.SubPackage {
	var gqlSubPackages []*gqlModel.SubPackage
	for _, s := range subPackages {
		gqlSubPackage := MapToGqlSubPackage(&s)
		gqlSubPackages = append(gqlSubPackages, gqlSubPackage)
	}
	return gqlSubPackages
}