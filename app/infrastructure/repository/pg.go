package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func SqlConnection(dns string) *sql.DB {
	db, err := sql.Open("postgres", dns)
	if err != nil {
		panic(err)
	}

	return db
}
