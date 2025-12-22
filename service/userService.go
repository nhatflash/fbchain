package service

import (
	"context"
	"time"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/enum"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/security"
)

type IUserService interface {
	GetCurrentUser(ctx context.Context) (*model.User, error)
	IsUserRoleTenant(u *model.User) bool
	GetListUser() ([]model.User, error)
	GetUserById(id int64) (*model.User, error)
	ChangeProfile(ctx context.Context, firstName *string, lastName *string, birthdate *string, gender *enum.Gender, phone *string, identity *string, address *string, postalCode *string, profileImage *string) (*model.User, error)
}

type UserService struct {
	UserRepo 		*repository.UserRepository
}


func NewUserService(ur *repository.UserRepository) IUserService {
	return &UserService{
		UserRepo: ur,
	}
}

func (us *UserService) GetCurrentUser(ctx context.Context) (*model.User, error) {
	var err error
	var claims *security.JwtAccessClaims
	claims, err = GetCurrentClaims(ctx)
	if err != nil {
		return nil, err
	}
	email := claims.Email
	var user *model.User
	user, err = us.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (*UserService) IsUserRoleTenant(u *model.User) bool {
	return u.Role == enum.ROLE_TENANT
}

func (us *UserService) GetListUser() ([]model.User, error) {
	users, err := us.UserRepo.ListAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}


func (us *UserService) GetUserById(id int64) (*model.User, error) {
	var err error
	var user *model.User
	user, err = us.UserRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}


func (us *UserService) ChangeProfile(ctx context.Context, firstName *string, lastName *string, birthdate *string, gender *enum.Gender, phone *string, identity *string, address *string, postalCode *string, profileImage *string) (*model.User, error) {
	var err error
	var claims *security.JwtAccessClaims
	claims, err = GetCurrentClaims(ctx)
	if err != nil {
		return nil, err
	}
	var u *model.User
	u, err = us.GetUserById(claims.UserId)
	if err != nil {
		return nil, err
	}
	var uFirstName *string
	var uLastName *string
	var uBirthdate *time.Time
	var uGender *enum.Gender
	var uPhone *string
	var uIdentity *string
	var uAddress *string
	var uPostalCode *string
	var uProfileImage *string

	uFirstName, uLastName, uBirthdate, uGender, uPhone, uIdentity, uAddress, uPostalCode, uProfileImage, err = GetDataForUserUpdate(firstName, lastName, birthdate, gender, phone, identity, address, postalCode, profileImage, u)
	if err != nil {
		return nil, err
	}

	err = ValidateChangeProfileRequest(phone, identity, uPhone, uIdentity, us.UserRepo)
	if err != nil {
		return nil, err
	}

	var updatedUser *model.User
	updatedUser, err = us.UserRepo.UpdateUser(ctx, u.Id, uFirstName, uLastName, uBirthdate, uGender, uPhone, uIdentity, uAddress, uPostalCode, uProfileImage)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}


func GetDataForUserUpdate(firstName *string, lastName *string, birthdate *string, gender *enum.Gender, phone *string, identity *string, address *string, postalCode *string, profileImage *string, u *model.User) (*string, *string, *time.Time, *enum.Gender, *string, *string, *string, *string, *string, error)  {
	if firstName == nil || *firstName == "" {
		firstName = &u.FirstName
	}
	if lastName == nil || *lastName == "" {
		lastName = &u.LastName
	}
	var bd *time.Time
	var err error

	if birthdate == nil || *birthdate == "" {
		bd = &u.Birthdate
	} else {
		bd, err = helper.ConvertToDate(*birthdate)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
	}
	if gender == nil {
		gender = &u.Gender
	}
	if phone == nil || *phone == "" {
		phone = u.Phone
	}
	if identity == nil || *identity == "" {
		identity = u.Identity
	}
	if address == nil || *address == "" {
		address = u.Address
	}
	if postalCode == nil || *postalCode == "" {
		postalCode = u.PostalCode
	}
	if profileImage == nil || *profileImage == "" {
		profileImage = u.ProfileImage
	}
	return firstName, lastName, bd, gender, phone, identity, address, postalCode, profileImage, err
}

func ValidateChangeProfileRequest(phone *string, identity *string, uPhone *string, uIdentity *string, ur *repository.UserRepository) error {
	if phone != nil && uPhone != nil && *phone != *uPhone && ur.CheckUserPhoneExists(*phone) {
		return appErr.BadRequestError("This phone is already in use.")
	}
	if identity != nil && uIdentity != nil && *identity != *uIdentity && ur.CheckUserIdentityExists(*identity) {
		return appErr.BadRequestError("This identity is already in use.")
	}
	return nil
}





