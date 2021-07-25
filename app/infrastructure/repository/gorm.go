package repository

import (
	"github.com/EdlanioJ/kbu-store/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GORMConnection(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	if cfg.Env != "test" {
		db, err = gorm.Open(postgres.Open(cfg.Dns), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.DnsTest), &gorm.Config{
			SkipDefaultTransaction: true,
		})
	}

	if err != nil {
		panic(err)
	}

	return db
}
