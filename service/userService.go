package service

import (
	"context"
	"github.com/nhatflash/fbchain/enum"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/security"
)

type IUserService interface {
	GetCurrentUser(ctx context.Context) (*model.User, error)
	IsUserRoleTenant(u *model.User) bool
	GetListUser() ([]model.User, error)
	GetUserById(id int64) (*model.User, error)
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


