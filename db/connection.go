package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type Config struct {
	Host         string
	Port         string
	PostgresUser string
	Password     string
	DbName       string
	SSLMode      string
}

var cfg Config

func Connect() (*sql.DB, error) {
	if cfg == (Config{}) {
		cfg = Config{
			PostgresUser: os.Getenv("DB_USERNAME"),
			Password:     os.Getenv("DB_PASSWORD"),
			DbName:       os.Getenv("DB_NAME"),
			Host:         os.Getenv("DB_HOST"),
			Port:         os.Getenv("DB_PORT"),
			SSLMode:      os.Getenv("SSL_MODE"),
		}
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.PostgresUser, cfg.Password, cfg.DbName, cfg.SSLMode,
	)

	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
		return db, err
	}

	return db, nil
}
