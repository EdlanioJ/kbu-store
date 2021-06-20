package utils

import (
	"github.com/EdlanioJ/kbu-store/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	var db *gorm.DB
	var err error

	config, err := config.LoadConfig("../../")
	if err != nil {
		panic(err)
	}

	if config.Env != "test" {
		db, err = gorm.Open(postgres.Open(config.Dns), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open(config.DnsTest), &gorm.Config{
			SkipDefaultTransaction: true,
		})
	}

	if err != nil {
		panic(err)
	}
	return db
}
