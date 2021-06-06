package utils

import (
	"github.com/EdlanioJ/kbu-store/app/config"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	var db *gorm.DB
	var err error

	config, err := config.LoadConfig("../../")

	if err != nil {
		logrus.Fatal("can not import config file", err)
	}

	if config.Env != "test" {
		db, err = gorm.Open(postgres.Open(config.PGDns), &gorm.Config{})
	} else {

		db, err = gorm.Open(sqlite.Open(config.PGDnsTest), &gorm.Config{
			SkipDefaultTransaction: true,
		})
	}

	if err != nil {
		logrus.Fatalf("Error connecting to database: %v", err)
	}

	if config.AutoMigration {
		_ = db.AutoMigrate(&dto.AccountDBModel{}, &dto.StoreDBModel{}, &dto.CategoryDBModel{})
	}

	return db
}
