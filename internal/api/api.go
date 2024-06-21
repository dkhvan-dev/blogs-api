package api

import (
	"blogs-api/internal/users/service"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	users service.Users
	DB    *sql.DB
}

func New(u service.Users, db *sql.DB) *API {
	return &API{
		users: u,
		DB:    db,
	}
}

func (a *API) AddRoutes(r *mux.Router) error {
	r = r.PathPrefix("/api/v1").Subrouter()

	r.HandleFunc("/auth", a.Login).Methods(http.MethodPost)

	r.HandleFunc("/users", a.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{login}", a.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/users", a.FindAllUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{login}", a.FindUserByLogin).Methods(http.MethodGet)
	r.HandleFunc("/users/{login}/password", a.ChangePassword).Methods(http.MethodPut)
	r.HandleFunc("/users/{login}", a.DeleteUser).Methods(http.MethodDelete)

	return nil
}
