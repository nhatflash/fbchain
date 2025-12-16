package service

import (
	"context"
	"database/sql"
	"github.com/nhatflash/fbchain/enum"
	appError "github.com/nhatflash/fbchain/error"
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
	Db *sql.DB
}


func NewUserService(db *sql.DB) IUserService {
	return &UserService{
		Db: db,
	}
}

func (u *UserService) GetCurrentUser(ctx context.Context) (*model.User, error) {
	var err error
	var claims *security.JwtAccessClaims
	claims, err = GetCurrentClaims(ctx)
	if err != nil {
		return nil, err
	}
	email := claims.Email
	var user *model.User
	user, err = repository.GetUserByEmail(email, u.Db)
	if err != nil {
		return nil, appError.NotFoundError("User not found")
	}
	return user, nil
}

func (*UserService) IsUserRoleTenant(u *model.User) bool {
	return *u.Role == enum.ROLE_TENANT
}

func (u *UserService) GetListUser() ([]model.User, error) {
	users, err := repository.ListAllUsers(u.Db)
	if err != nil {
		return nil, err
	}
	return users, nil
}


func (u *UserService) GetUserById(id int64) (*model.User, error) {
	var err error
	var user *model.User
	user, err = repository.GetUserById(u.Db, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}


