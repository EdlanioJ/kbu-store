package repository

import (
	"fmt"

	"github.com/EdlanioJ/kbu-store/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GORMConnection(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
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

		db, err = gorm.Open(postgres.Open(dns))
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.DBTest), &gorm.Config{
			SkipDefaultTransaction: true,
		})
	}

	if err != nil {
		panic(err)
	}

	return db
}
