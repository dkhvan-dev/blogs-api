package api

import (
	"blogs-api/internal/middlewares"
	"blogs-api/internal/users/service"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

const ADMIN_ROLE = "ADMIN"

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
	appRouter := r.PathPrefix("/api/v1").Subrouter()

	addUserRoutes(appRouter, a)
	addAdminRoutes(appRouter, a)

	r.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("/blogs-api/cmd/docs/"))))

	return nil
}

func addUserRoutes(r *mux.Router, a *API) {
	r.HandleFunc("/auth", a.Login).Methods(http.MethodPost)
	r.Handle("/logout", middlewares.AuthMiddleware("", http.HandlerFunc(a.Logout))).Methods(http.MethodPost)

	r.HandleFunc("/users", a.CreateUser).Methods(http.MethodPost)
	r.Handle("/users", middlewares.AuthMiddleware("", http.HandlerFunc(a.UpdateUser))).Methods(http.MethodPut)
	r.Handle("/users/password", middlewares.AuthMiddleware("", http.HandlerFunc(a.ChangePassword))).Methods(http.MethodPut)
}

func addAdminRoutes(r *mux.Router, a *API) {
	r.Handle("/admin/users", middlewares.AuthMiddleware(ADMIN_ROLE, http.HandlerFunc(a.FindAllUsers))).Methods(http.MethodGet)
	r.Handle("/admin/users/{login}", middlewares.AuthMiddleware(ADMIN_ROLE, http.HandlerFunc(a.FindUserByLogin))).Methods(http.MethodGet)
	r.Handle("/admin/users/{login}", middlewares.AuthMiddleware(ADMIN_ROLE, http.HandlerFunc(a.DeleteUser))).Methods(http.MethodDelete)
}
