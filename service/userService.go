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
	"github.com/nhatflash/fbchain/client"
)

type IUserService interface {
	FindCurrentUser(ctx context.Context) (*model.User, error)
	IsUserRoleTenant(u *model.User) bool
	FindAllUsers(ctx context.Context) ([]model.User, error)
	FindUserById(ctx context.Context, id int64) (*model.User, error)
	HandleChangeUserProfile(ctx context.Context, user *model.User, req *client.UpdateProfileRequest) (*model.User, error)
}

type UserService struct {
	UserRepo 		*repository.UserRepository
}


func NewUserService(ur *repository.UserRepository) IUserService {
	return &UserService{
		UserRepo: ur,
	}
}

func (us *UserService) FindCurrentUser(ctx context.Context) (*model.User, error) {
	var err error
	var claims *security.JwtAccessClaims
	claims, err = GetCurrentClaims(ctx)
	if err != nil {
		return nil, err
	}
	email := claims.Email
	var user *model.User
	user, err = us.UserRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (*UserService) IsUserRoleTenant(u *model.User) bool {
	return u.Role == enum.ROLE_TENANT
}

func (us *UserService) FindAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := us.UserRepo.FindAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}


func (us *UserService) FindUserById(ctx context.Context, id int64) (*model.User, error) {
	var err error
	var user *model.User
	user, err = us.UserRepo.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}


func (us *UserService) HandleChangeUserProfile(ctx context.Context, user *model.User, req *client.UpdateProfileRequest) (*model.User, error) {
	firstName, lastName, birthdate, gender, phone, identity, address, postalCode, profileImage, err := getDataForUserUpdate(req.FirstName, 
						 req.LastName, 
						 req.Birthdate, 
						 req.Gender, 
						 req.Phone, 
						 req.Identity, 
						 req.Address, 
						 req.PostalCode, 
						 req.ProfileImage, 
						 user)
	if err != nil {
		return nil, err
	}

	if err := validateChangeProfileRequest(ctx, req.Phone, req.Identity, phone, identity, us.UserRepo); err != nil {
		return nil, err
	}
	
	updatedUser, err := us.UserRepo.UpdateUserInfo(ctx, user.Id, firstName, lastName, birthdate, gender, phone, identity, address, postalCode, profileImage)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}


func getDataForUserUpdate(firstName *string, lastName *string, birthdate *string, gender *enum.Gender, phone *string, identity *string, address *string, postalCode *string, profileImage *string, u *model.User) (*string, *string, *time.Time, *enum.Gender, *string, *string, *string, *string, *string, error)  {
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

func validateChangeProfileRequest(ctx context.Context, reqPhone *string, reqIdentity *string, phone *string, identity *string, ur *repository.UserRepository) error {
	if reqPhone != nil && phone != nil && *reqPhone != *phone {
		exist, err := ur.CheckUserPhoneExists(ctx, *reqPhone)
		if err != nil {
			return err
		}
		if exist {
			return appErr.BadRequestError("This phone is already in use.")
		}
	}
	if reqIdentity != nil && identity != nil && *reqIdentity != *identity {
		exist, err := ur.CheckUserIdentityExists(ctx, *reqIdentity)
		if err != nil {
			return err
		}
		if exist {
			return appErr.BadRequestError("This identity is already in use.")
		}
	}
	return nil
}





