package service

import (
	"blogs-api/internal/users/model"
)

type UserService interface {
	FindAll(page, size int, search string) []model.UserView
	FindById(login string) (*model.UserView, error)
	Create(input model.UserCreate) error
	Update(input model.UserUpdate) error
	ChangePassword(login string, input model.UserChangePassword) error
	Delete(login string) error
}

type Users struct {
	Service UserService
}

func New(s UserService) *Users {
	return &Users{
		Service: s,
	}
}
