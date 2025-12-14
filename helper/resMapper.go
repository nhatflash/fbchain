package helper

import (
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/model"
)

func MapToUserResponse(u *model.User) *client.UserResponse {
	phone, identity, postalCode, address, profileImage := getUserDataIfSqlNullString(u)
	userRes := client.UserResponse{
		Id:           u.Id,
		Email:        u.Email,
		Role:         u.Role,
		Phone:        phone,
		Identity:     identity,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Gender:       u.Gender,
		Birthdate:    u.Birthdate,
		PostalCode:   postalCode,
		Address:      address,
		ProfileImage: profileImage,
		Status:       u.Status,
	}
	return &userRes
}

func MapToTenantResponse(u *model.User, t *model.Tenant) *client.TenantResponse {
	phone, identity, postalCode, address, profileImage := getUserDataIfSqlNullString(u)
	description, notes := getTenantDataIfSqlNullString(t)

	tenantRes := client.TenantResponse{
		UserId:       u.Id,
		Email:        u.Email,
		Phone:        phone,
		Identity:     identity,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Gender:       u.Gender,
		Birthdate:    u.Birthdate,
		PostalCode:   postalCode,
		Address:      address,
		ProfileImage: profileImage,
		Code:         t.Code,
		Description:  description,
		Type:         t.Type,
		Notes:        notes,
		Status:       u.Status,
	}
	return &tenantRes
}

func MapToSubscriptionResponse(s *model.Subscription) *client.SubscriptionResponse {
	description, image := getSubscriptionDataIfSqlNullString(s)
	subscriptionRes := client.SubscriptionResponse{
		Id:            s.Id,
		Name:          s.Name,
		Description:   description,
		DurationMonth: s.DurationMonth,
		Price:         s.Price,
		IsActive:      s.IsActive,
		Image:         image,
	}
	return &subscriptionRes
}

func MapToRestaurantResponse(r *model.Restaurant) *client.RestaurantResponse {
	description, contactEmail, contactPhone, notes := getRestaurantDataIfSqlNullString(r)

	var images []string
	for i := range r.Images {
		image := r.Images[i]
		images = append(images, image.Image)
	}
	restaurantRes := client.RestaurantResponse{
		Id:             r.Id,
		TenantId:       r.TenantId,
		Name:           r.Name,
		Location:       r.Location,
		Description:    description,
		ContactEmail:   contactEmail,
		ContactPhone:   contactPhone,
		PostalCode:     r.PostalCode,
		Type:           r.Type,
		AvgRating:      r.AvgRating,
		Notes:          notes,
		SubscriptionId: r.SubscriptionId,
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

func getUserDataIfSqlNullString(u *model.User) (string, string, string, string, string) {
	phone := ""
	identity := ""
	postalCode := ""
	address := ""
	profileImage := ""
	if u.Phone.Valid {
		phone = u.Phone.String
	}
	if u.Identity.Valid {
		identity = u.Identity.String
	}
	if u.Address.Valid {
		address = u.Address.String
	}
	if u.ProfileImage.Valid {
		profileImage = u.ProfileImage.String
	}
	if u.PostalCode.Valid {
		postalCode = u.PostalCode.String
	}
	return phone, identity, postalCode, address, profileImage
}

func getTenantDataIfSqlNullString(t *model.Tenant) (string, string) {
	description := ""
	notes := ""
	if t.Description.Valid {
		description = t.Description.String
	}
	if t.Notes.Valid {
		notes = t.Notes.String
	}
	return description, notes
}

func getSubscriptionDataIfSqlNullString(s *model.Subscription) (string, string) {
	description := ""
	image := ""
	if s.Description.Valid {
		description = s.Description.String
	}
	if s.Image.Valid {
		image = s.Image.String
	}
	return description, image
}

func getRestaurantDataIfSqlNullString(r *model.Restaurant) (string, string, string, string) {
	description := ""
	contactEmail := ""
	contactPhone := ""
	notes := ""
	if r.Description.Valid {
		description = r.Description.String
	}
	if r.ContactEmail.Valid {
		contactEmail = r.ContactEmail.String
	}
	if r.ContactPhone.Valid {
		contactPhone = r.ContactPhone.String
	}
	if r.Notes.Valid {
		notes = r.Notes.String
	}
	return description, contactEmail, contactPhone, notes
}
