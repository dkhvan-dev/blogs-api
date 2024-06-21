package store

import (
	"blogs-api/internal/users/model"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserStore struct {
	db *sql.DB
}

func New(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) FindAll(page, size int, search string) []model.UserView {
	var response []model.UserView
	rows, err := s.db.Query("select * from t_users limit ? offset ?", size, (page+1)*size)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var user model.User
		err = rows.Scan(&user)

		if err != nil {
			panic(err)
		}

		response = append(response, *user.ToView())
	}

	return response
}

func (s *UserStore) FindById(login string) (*model.UserView, error) {
	row := s.db.QueryRow("select * from t_users where login = ? limit 1", login)
	var user model.User
	err := row.Scan(&user)

	if err != nil {
		return nil, errors.New("User not found!")
	}

	return user.ToView(), nil
}

func (s *UserStore) Create(input model.UserCreate) error {
	err := validateCreationUser(input)

	if err != nil {
		return err
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	_, err = s.db.Query("insert into t_users (login, password, email, first_name, middle_name, last_name, birthdate) values (?, ?, ?, ?, ?, ?, ?)", input.Login, string(hashPassword), input.Email, input.FirstName, input.MiddleName, input.LastName, input.Birthdate.Format("2006-02-15"))

	if err != nil {
		return errors.New("User creation failed!")
	}

	return nil
}

func (s *UserStore) Update(input model.UserUpdate) error {
	return nil
}

func (s *UserStore) ChangePassword(login string, input model.UserChangePassword) error {
	return nil
}

func (s *UserStore) Delete(login string) error {
	return nil

}

func validateCreationUser(input model.UserCreate) error {
	currentDate, err := time.Parse(time.DateOnly, time.Now().String())

	if err != nil {
		return errors.New("Invalid birthdate!")
	}

	if input.Password != input.RePassword {
		return errors.New("Passwords didn't match!")
	}

	if input.Birthdate.After(currentDate) {
		return errors.New("Invalid birthdate!")
	}

	return nil
}
