package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host         = "localhost"
	port         = 5432
	postgresUser = "postgres"
	password     = "example"
	dbname       = "plusz-db"
)

func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, postgresUser, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}

	return db, nil
}
