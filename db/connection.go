package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host         = "ec2-54-195-190-73.eu-west-1.compute.amazonaws.com"
	port         = 5432
	postgresUser = "ufqgndot5qvr7m"
	password     = "p6c7fe848c130f3b09afb2c6a514aa54e329c775e8e8bda95c13dffba72dbc4ea"
	dbname       = "db98h2g5gecn0e"
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
