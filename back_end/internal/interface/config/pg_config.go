package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgreSQLConnection() (*sqlx.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		panic(fmt.Errorf("error connecting to: %w", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("error ping db: %w", err))
	}

	return db, nil
}
