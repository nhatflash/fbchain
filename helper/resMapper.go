package helper

import (
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/client"
)


func MapToUserResponse(u *model.User) *client.UserResponse {
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
	userRes := client.UserResponse{
		Id: u.Id,
		Email: u.Email,
		Role: u.Role,
		Phone: phone,
		Identity: identity,
		FirstName: u.FirstName,
		LastName: u.LastName,
		Gender: u.Gender,
		Birthdate: u.Birthdate,
		PostalCode: postalCode,
		Address: address,
		ProfileImage: profileImage,
		Status: u.Status,
	}
	return &userRes
}