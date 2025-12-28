package graph

import (
	"strconv"

	gqlModel "github.com/nhatflash/fbchain/graph/model"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/scalar"
)

func MapToGqlModelUser(u *model.User) *gqlModel.User {
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

func MapToGqlModelTenant(t *model.Tenant) *gqlModel.Tenant {
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

func MapToGqlRestaurant(r *model.Restaurant) *gqlModel.Restaurant {
	idStr := strconv.FormatInt(r.TenantId, 10)
	avgRating64,_ := r.AvgRating.Float64()
	return &gqlModel.Restaurant{
		ID: strconv.FormatInt(r.Id, 10),
		Name: r.Name,
		Location: &r.Location,
		Description: r.Description,
		ContactEmail: r.ContactEmail,
		ContactPhone: r.ContactPhone,
		PostalCode: &r.PostalCode,
		Type: &r.Type,
		AvgRating: &avgRating64,
		IsActive: &r.IsActive,
		Notes: r.Notes,
		TenantID: &idStr,
	}
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


func MapToGqlRestaurantItem(item *model.RestaurantItem) *gqlModel.RestaurantItem {
	idStr := strconv.FormatInt(item.RestaurantId, 10)
	price := item.Price.String()
	return &gqlModel.RestaurantItem{
		ID: item.Id,
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