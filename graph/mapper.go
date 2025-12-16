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
		Gender:       (*string)(&u.Gender),
		Birthdate:    (*scalar.CustomDate)(&u.Birthdate),
		PostalCode:   u.PostalCode,
		Address:      u.Address,
		ProfileImage: u.ProfileImage,
	}
}