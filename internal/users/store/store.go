package store

import (
	"blogs-api/internal/core/errors"
	"blogs-api/internal/users/model"
	"blogs-api/internal/utils"
	"database/sql"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"math"
	"strconv"
	"time"
)

var JwtKey = []byte(utils.GetEnv("JWT_SECRET"))
var JwtLifespanMinutes = []byte(utils.GetEnv("JWT_LIFESPAN_MINUTES"))

const (
	ErrCreation           = errors.Error("User creation failed")
	ErrUpdating           = errors.Error("User updating failed")
	ErrParsingBirthdate   = errors.Error("Parsing birthdate failed")
	ErrParsingCurrentDate = errors.Error("Parsing current date failed")
	ErrPasswordsMatch     = errors.Error("Passwords didn't match")
	ErrNewPasswordsMatch  = errors.Error("New passwords does not match")
	ErrInvalidBirthdate   = errors.Error("Invalid birthdate!")
	ErrUserAlreadyExists  = errors.Error("User already exists")
	ErrUserNotFound       = errors.Error("User not found")
	ErrWrongOldPassword   = errors.Error("Wrong old password")
	ErrPasswordEncryption = errors.Error("Password encryption failed")
	ErrChangePassword     = errors.Error("Change password failed")
	ErrUserDelete         = errors.Error("Deleting user failed")
	ErrLoginData          = errors.Error("Invalid login or password")
)

type UserStore struct {
	db *sql.DB
}

type JwtToken struct {
	AccessToken            string
	IssuedAt               time.Time
	ExpiresAt              time.Time
	RefreshToken           string
	RefreshTokenExpireTime time.Time
}

func New(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) Login(input model.UserLoginRequest) (*map[string]string, error) {
	entity, err := s.GetEntityByLogin(input.Login)
	if err != nil {
		return nil, ErrUserNotFound.Wrap(errors.ErrNotFound)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(input.Password)); err != nil {
		return nil, ErrLoginData.Wrap(errors.ErrBadRequest)
	}

	accessToken, refreshToken, err := generateToken(input.Login)

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return &map[string]string{
		"access_token":             accessToken,
		"access_token_expire_time": time.Now().Add(strconv.ParseInt(utils.GetEnv(JwtLifespanMinutes)) * time.Minute),
		"refresh_token":            refreshToken,
	}, nil
}

func (s *UserStore) FindAll(page, size int, search string) []model.UserView {
	if size == 0 {
		size = math.MaxInt
	}

	response := []model.UserView{}
	rows, err := s.db.Query("select * from t_users offset $1 limit $2", page, size)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var user model.User
		var role model.Role
		var roleCode string
		err = rows.Scan(&user.Login, &user.Password, &user.Email, &user.FirstName, &user.MiddleName, &user.LastName, &user.Birthdate, &roleCode)

		if err != nil {
			panic(err)
		}

		roleRow := s.db.QueryRow("select * from t_roles where code = $1", roleCode)
		err = roleRow.Scan(&role.Code, &role.Name)
		if err != nil {
			panic(err)
		}

		user.Role = role

		response = append(response, *user.ToView())
	}

	return response
}

func (s *UserStore) FindByLogin(login string) (*model.UserView, error) {
	user, err := s.GetEntityByLogin(login)

	if err != nil {
		return nil, ErrUserNotFound.Wrap(errors.ErrNotFound)
	}

	return user.ToView(), nil
}

func (s *UserStore) GetEntityByLogin(login string) (*model.User, error) {
	var user model.User
	var role model.Role
	var roleCode string

	row := s.db.QueryRow("select * from t_users where login = $1", login)
	err := row.Scan(&user.Login, &user.Password, &user.Email, &user.FirstName, &user.MiddleName, &user.LastName, &user.Birthdate, &roleCode)

	if err != nil {
		return nil, ErrUserNotFound.Wrap(errors.ErrNotFound)
	}

	roleRow := s.db.QueryRow("select * from t_roles where code = $1", roleCode)
	err = roleRow.Scan(&role.Code, &role.Name)
	if err != nil {
		panic(err)
	}

	user.Role = role

	return &user, nil
}

func (s *UserStore) Create(input model.UserCreate) error {
	exists := s.ExistsByLogin(input.Login)
	if exists {
		return ErrUserAlreadyExists.Wrap(errors.ErrBadRequest)
	}

	err := validateCreationUser(input)
	if err != nil {
		return err
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	_, err = s.db.Exec("insert into t_users (login, password, email, first_name, middle_name, last_name, birthdate) values ($1, $2, $3, $4, $5, $6, $7)", input.Login, string(hashPassword), input.Email, input.FirstName, input.MiddleName, input.LastName, input.Birthdate)

	if err != nil {
		return ErrCreation.Wrap(errors.ErrUnknown)
	}

	return nil
}

func (s *UserStore) Update(login string, input model.UserUpdate) error {
	exists := s.ExistsByLogin(login)

	if !exists {
		return ErrUserNotFound.Wrap(errors.ErrNotFound)
	}

	_, err := s.db.Exec("update t_users set firstName = $2, middleName = $3, lastName = $4, birthdate = $5 where login = $1", login, input.FirstName, input.MiddleName, input.LastName, input.Birthdate)
	if err != nil {
		return ErrUpdating.Wrap(errors.ErrUnknown)
	}

	return nil
}

func (s *UserStore) ChangePassword(login string, input model.UserChangePassword) error {
	exists := s.ExistsByLogin(login)
	if !exists {
		return ErrUserNotFound.Wrap(errors.ErrNotFound)
	}

	oldPassword := ""
	row := s.db.QueryRow("select password from t_users where login = $1", login)
	row.Scan(&oldPassword)

	if err := bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(input.OldPassword)); err != nil {
		return ErrWrongOldPassword.Wrap(errors.ErrBadRequest)
	}

	if input.NewPassword != input.RePassword {
		return ErrNewPasswordsMatch.Wrap(errors.ErrBadRequest)
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return ErrPasswordEncryption.Wrap(errors.ErrUnknown)
	}

	_, err = s.db.Exec("update t_users set password = $1 where login = $2", newPassword, login)
	if err != nil {
		return ErrChangePassword.Wrap(errors.ErrUnknown)
	}

	return nil
}

func (s *UserStore) Delete(login string) error {
	if !s.ExistsByLogin(login) {
		return ErrUserNotFound.Wrap(errors.ErrNotFound)
	}

	_, err := s.db.Exec("delete from t_users where login = $1", login)
	if err != nil {
		return ErrUserDelete.Wrap(errors.ErrUnknown)
	}

	return nil
}

func (s *UserStore) ExistsByLogin(login string) bool {
	userExists := false
	row := s.db.QueryRow("select exists(select 1 from t_users where login = $1)", login)
	row.Scan(&userExists)

	return userExists
}

func validateCreationUser(input model.UserCreate) error {
	birthdate, err := time.Parse("2006-01-02", input.Birthdate)
	if err != nil {
		return ErrParsingBirthdate.Wrap(errors.ErrValidation)
	}

	currentTime := time.Now().Format("2006-02-15")
	currentDate, err := time.Parse("2006-02-15", currentTime)

	if err != nil {
		return ErrParsingCurrentDate.Wrap(errors.ErrValidation)
	}

	if input.Password != input.RePassword {
		return ErrPasswordsMatch.Wrap(errors.ErrBadRequest)
	}

	if birthdate.After(currentDate) {
		return ErrInvalidBirthdate.Wrap(errors.ErrBadRequest)
	}

	return nil
}

func generateToken(username string) (string, string, error) {
	expirationTime := time.Now().Add(15 * time.Minute).Unix()

	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime,
		Issuer:    username,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString(JwtKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = 1
	refreshTokenClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	refreshTokenString, err := refreshToken.SignedString(JwtKey)

	return accessTokenString, refreshTokenString, nil
}
