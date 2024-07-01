// @title User Service API
// @version 1.0
// @description User Service.

// @host localhost:8080
// @BasePath /api/v1

package main

import (
	_ "blogs-api/cmd/docs"
	"blogs-api/internal/api"
	"blogs-api/internal/db"
	"blogs-api/internal/users/service"
	"blogs-api/internal/users/store"
	"blogs-api/internal/utils"
	"blogs-api/pkg"
	"database/sql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"net/http"
)

const SERVER_PORT = "SERVER_PORT"

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	port := ":" + utils.GetEnv(SERVER_PORT)
	userStore := store.New(db.DB)
	userService := service.New(userStore)
	srv := api.New(*userService, db.DB)
	r := mux.NewRouter()
	srv.AddRoutes(r)

	scheduler := cron.New()

	_, err := scheduler.AddFunc("30 * * * *", func() {
		pkg.Logger.Info("Starting schedule delete tokens")
		deleteTokens(db.DB)
		pkg.Logger.Info("Finishing schedule delete tokens")
	})

	if err != nil {
		panic(err)
	}

	scheduler.Start()

	pkg.Logger.Fatal(http.ListenAndServe(port, r).Error())
}

func deleteTokens(db *sql.DB) {
	if _, err := db.Exec("delete from t_tokens where expiration_deadline >= now()"); err != nil {
		panic(err)
	}
}
