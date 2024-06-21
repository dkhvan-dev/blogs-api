package main

import (
	"blogs-api/internal/api"
	"blogs-api/internal/db"
	"blogs-api/internal/users/service"
	"blogs-api/internal/users/store"
	"blogs-api/internal/utils"
	"blogs-api/pkg"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
)

type Server struct {
	api    *api.API
	router *mux.Router
}

func main() {
	pkg.InitLogger()
	pkg.Logger.Info("Starting server")

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// init db
	db.InitDB()

	// init server
	port := ":" + utils.GetEnv("SERVER_PORT")
	userStore := store.New(db.DB)
	userService := service.New(userStore)
	srv := api.New(*userService, db.DB)

	// start server
	start(srv, port)
}

func start(srv *api.API, port string) {
	r := mux.NewRouter()
	srv.AddRoutes(r)

	pkg.Logger.Fatal(http.ListenAndServe(port, r).Error())
}
