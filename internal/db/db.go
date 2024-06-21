package db

import (
	"blogs-api/internal/utils"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"os"
)

var DB *sql.DB

func InitDB() {
	dbHost := utils.GetEnv("DB_HOST")
	dbName := utils.GetEnv("DB_NAME")
	dbUser := utils.GetEnv("DB_USER")
	dbPass := utils.GetEnv("DB_PASS")
	dbPort := utils.GetEnv("DB_PORT")
	sslmode := utils.GetEnv("SSLMODE")
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, sslmode)

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}

	if db.Ping() != nil {
		panic(err)
	}

	DB = db

	err = migrateSql(db)
	if err != nil {
		panic(err)
	}
}

func migrateSql(db *sql.DB) error {
	migrationsPath := "file://internal/migrations"
	dbName := os.Getenv("DB_NAME")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(migrationsPath, dbName, driver)

	if err != nil {
		return err
	}

	defer migrator.Close()

	if err = migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
