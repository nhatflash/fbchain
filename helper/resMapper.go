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
	return &client.UserResponse{
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
}

func MapToTenantResponse(u *model.User, t *model.Tenant) *client.TenantResponse {
	return &client.TenantResponse{
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
}

func MapToSubPackageResponse(s *model.SubPackage) *client.SubPackageResponse {
	return &client.SubPackageResponse{
		Id:            s.Id,
		Name:          s.Name,
		Description:   s.Description,
		DurationMonth: s.DurationMonth,
		Price:         s.Price,
		IsActive:      s.IsActive,
		Image:         s.Image,
	}
}

func MapToRestaurantResponse(r *model.Restaurant, rImgs []model.RestaurantImage) *client.RestaurantResponse {
	var images []string
	for i := range rImgs {
		image := rImgs[i]
		images = append(images, image.Image)
	}
	return &client.RestaurantResponse{
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
}

func MapToOrderResponse(o *model.Order) *client.OrderResponse {
	return &client.OrderResponse{
		Id:           o.Id,
		TenantId:     o.TenantId,
		RestaurantId: o.RestaurantId,
		OrderDate:    o.OrderDate,
		Status:       o.Status,
		Amount:       o.Amount,
	}
}


func MapToRestaurantItemResponse(i *model.RestaurantItem) *client.RestaurantItemResponse {
	price := i.Price.String()
	itemId := i.Id.Hex()
	return &client.RestaurantItemResponse{
		Id: 		itemId,
		Name: 		i.Name,
		Description: i.Description,
		Price: 		price,
		Type: 		i.Type,
		Status: 	i.Status,
		Image: 		i.Image,
		Notes: 		i.Notes,
		RestaurantId: i.RestaurantId,
	}
}


func MapToRestaurantItemsResponse(items []model.RestaurantItem) []client.RestaurantItemResponse {
	var itemsRes []client.RestaurantItemResponse
	for _, i := range items {
		itemRes := MapToRestaurantItemResponse(&i)
		itemsRes = append(itemsRes, *itemRes)
	}
	return itemsRes
}


func MapToRestaurantTableResponse(t *model.RestaurantTable) *client.RestaurantTableResponse {
	return &client.RestaurantTableResponse{
		Id: t.Id,
		RestaurantId: t.RestaurantId,
		Label: t.Label,
		IsActive: t.IsActive,
		Notes: t.Notes,
	}
}


func MapToRestaurantOrderResponse(o *model.RestaurantOrder) *client.RestaurantOrderResponse {
	var itemsRes []client.RestaurantOrderItemResponse
	for _, i := range o.Items {
		item := MapToRestaurantOrderItemResponse(&i)
		itemsRes = append(itemsRes, *item)
	}
	return &client.RestaurantOrderResponse{
		Id: o.Id,
		RestaurantId: o.RestaurantId,
		TableId: o.TableId,
		Amount: o.Amount,
		Status: o.Status,
		Notes: o.Notes,
		Items: itemsRes,
	}
}


func MapToRestaurantOrderItemResponse(i *model.RestaurantOrderItem) *client.RestaurantOrderItemResponse {
	return &client.RestaurantOrderItemResponse{
		Id: i.Id,
		ROrderId: i.ROrderId,
		ItemId: i.ItemId,
		Quantity: i.Quantity,
		Total: i.Total,
	}
}



