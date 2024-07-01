package service

import (
	"blogs-api/internal/core/errors"
	"blogs-api/internal/users/model"
)

type UserService interface {
	Login(input model.UserLoginRequest) (map[string]any, *errors.Error)
	FindAll(page, size int, search string) []model.UserView
	FindByLogin(login string) (*model.UserView, *errors.Error)
	GetEntityByLogin(login string) (*model.User, *errors.Error)
	Create(input model.UserCreate) *errors.Error
	Update(input model.UserUpdate) *errors.Error
	ChangePassword(input model.UserChangePassword) *errors.Error
	Delete(login string) *errors.Error
	ExistsByLogin(login string) bool
	Logout(login string) *errors.Error
}

type Users struct {
	Service UserService
}

func New(s UserService) *Users {
	return &Users{
		Service: s,
	}
}
