package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const driverName = "postgres"

var Dbx *sqlx.DB

func Connect() error {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort,
	)

	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return err
	}

	fmt.Println("Successfully connected to database")
	Dbx = db
	return nil
}

func Close() error {
	return Dbx.Close()
}
