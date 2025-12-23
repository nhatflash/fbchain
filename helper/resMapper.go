package helper

import (
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/model"
	"time"
)

func MapToSignInResponse(accessToken string, refreshToken string) *client.SignInResponse {
	signInRes := client.SignInResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		LastLogin: time.Now(),
	}
	return &signInRes
}


func MapToUserResponse(u *model.User) *client.UserResponse {
	userRes := client.UserResponse{
		Id:           u.Id,
		Email:        u.Email,
		Role:         u.Role,
		Phone:        u.Phone,
		Identity:     u.Identity,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Gender:       u.Gender,
		Birthdate:    u.Birthdate,
		PostalCode:   u.PostalCode,
		Address:      u.Address,
		ProfileImage: u.ProfileImage,
		Status:       u.Status,
	}
	return &userRes
}

func MapToTenantResponse(u *model.User, t *model.Tenant) *client.TenantResponse {
	tenantRes := client.TenantResponse{
		UserId:       u.Id,
		Email:        u.Email,
		Phone:        u.Phone,
		Identity:     u.Identity,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Gender:       u.Gender,
		Birthdate:    u.Birthdate,
		PostalCode:   u.PostalCode,
		Address:      u.Address,
		ProfileImage: u.ProfileImage,
		Code:         t.Code,
		Description:  t.Description,
		Type:         t.Type,
		Notes:        t.Notes,
		Status:       u.Status,
	}
	return &tenantRes
}

func MapToSubPackageResponse(s *model.SubPackage) *client.SubPackageResponse {
	subscriptionRes := client.SubPackageResponse{
		Id:            s.Id,
		Name:          s.Name,
		Description:   s.Description,
		DurationMonth: s.DurationMonth,
		Price:         s.Price,
		IsActive:      s.IsActive,
		Image:         s.Image,
	}
	return &subscriptionRes
}

func MapToRestaurantResponse(r *model.Restaurant, rImgs []model.RestaurantImage) *client.RestaurantResponse {
	var images []string
	for i := range rImgs {
		image := rImgs[i]
		images = append(images, image.Image)
	}
	restaurantRes := client.RestaurantResponse{
		Id:             r.Id,
		TenantId:       r.TenantId,
		Name:           r.Name,
		Location:       r.Location,
		Description:    r.Description,
		ContactEmail:   r.ContactEmail,
		ContactPhone:   r.ContactPhone,
		PostalCode:     r.PostalCode,
		Type:           r.Type,
		AvgRating:      r.AvgRating,
		Notes:          r.Notes,
		SubPackageId:   r.SubPackageId,
		Images:         images,
	}
	return &restaurantRes
}

func MapToOrderResponse(o *model.Order) *client.OrderResponse {
	orderRes := client.OrderResponse{
		Id:           o.Id,
		TenantId:     o.TenantId,
		RestaurantId: o.RestaurantId,
		OrderDate:    o.OrderDate,
		Status:       o.Status,
		Amount:       o.Amount,
	}
	return &orderRes
}


func MapToRestaurantItemResponse(i *model.RestaurantItem) *client.RestaurantItemResponse {
	itemRes := client.RestaurantItemResponse{
		Id: 		i.Id,
		Name: 		i.Name,
		Description: i.Description,
		Price: 		i.Price,
		Type: 		i.Type,
		Status: 	i.Status,
		Image: 		i.Image,
		Notes: 		i.Notes,
		RestaurantId: i.RestaurantId,
	}
	return &itemRes
}


