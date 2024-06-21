package main

import (
	"blogs-api/internal/api"
	"blogs-api/internal/db"
	"blogs-api/internal/users/service"
	"blogs-api/internal/users/store"
	"blogs-api/internal/utils"
	"blogs-api/pkg"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	start()
}

func start() {
	// init logger
	pkg.InitLogger()
	pkg.Logger.Info("Starting server")

	// load environments
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
	r := mux.NewRouter()
	srv.AddRoutes(r)

	pkg.Logger.Fatal(http.ListenAndServe(port, r).Error())
}
