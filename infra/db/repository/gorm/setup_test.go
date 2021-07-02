package gorm_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dbMock() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return gormDB, mock
}

func getAccount() *domain.Account {
	account, _ := domain.NewAccount(20.75)
	return account
}

func getCategory() *domain.Category {
	category, _ := domain.NewCategory("store type 001")
	return category
}
