package repository

import (
	"database/sql"
	"fmt"

	"github.com/EdlanioJ/kbu-store/app/config"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func SqlConnection(cfg *config.Config) *sql.DB {
	var db *sql.DB
	var err error

	if cfg.Env != "test" {
		dns := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
			cfg.PG.Host,
			cfg.PG.Port,
			cfg.PG.User,
			cfg.PG.Password,
			cfg.PG.Name,
			cfg.PG.SslMode,
		)
		db, err = sql.Open("postgres", dns)
	} else {
		db, err = sql.Open("sqlite", cfg.DBTest)
	}
	if err != nil {
		panic(err)
	}

	return db
}
