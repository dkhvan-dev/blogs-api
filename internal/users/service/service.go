package service

import (
	"blogs-api/internal/users/model"
)

type UserService interface {
	Login(input model.UserLoginRequest) (*map[string]string, error)
	FindAll(page, size int, search string) []model.UserView
	FindByLogin(login string) (*model.UserView, error)
	GetEntityByLogin(login string) (*model.User, error)
	Create(input model.UserCreate) error
	Update(login string, input model.UserUpdate) error
	ChangePassword(login string, input model.UserChangePassword) error
	Delete(login string) error
	ExistsByLogin(login string) bool
}

type Users struct {
	Service UserService
}

func New(s UserService) *Users {
	return &Users{
		Service: s,
	}
}
