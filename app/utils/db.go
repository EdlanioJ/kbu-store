package utils

import (
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB(env, dns, migration string) *gorm.DB {
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
		logrus.Fatalf("Error connecting to database: %v", err)
	}

	if migration == "true" {
		_ = db.AutoMigrate(&dto.AccountDBModel{}, &dto.StoreDBModel{}, &dto.CategoryDBModel{})
	}

	return db
}
