package repository

import (
	"database/sql"

	"github.com/EdlanioJ/kbu-store/app/config"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func SqlConnection(cfg *config.Config) *sql.DB {
	var db *sql.DB
	var err error

	if cfg.Env != "test" {
		db, err = sql.Open("postgres", cfg.Dns)
	} else {
		db, err = sql.Open("sqlite", cfg.DnsTest)
	}
	if err != nil {
		panic(err)
	}

	return db
}
