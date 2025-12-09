package helper

import (
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/client"
)


func MapToUserResponse(u *model.User) *client.UserResponse {
	userRes := client.UserResponse{
		Id: u.Id,
		Email: u.Email,
		Role: u.Role,
		Phone: u.Phone,
		Identity: u.Identity,
		FirstName: u.FirstName,
		LastName: u.LastName,
		Gender: u.Gender,
		Birthdate: u.Birthdate,
		PostalCode: u.PostalCode,
		Address: u.Address,
		ProfileImage: u.ProfileImage,
		Status: u.Status,
	}
	return &userRes
}