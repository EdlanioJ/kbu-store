package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GORMConnection(dns string, env string) *gorm.DB {
	var db *gorm.DB
	var err error

	if env != "test" {
		db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open(dns), &gorm.Config{
			SkipDefaultTransaction: true,
		})
	}

	if err != nil {
		panic(err)
	}

	return db
}
