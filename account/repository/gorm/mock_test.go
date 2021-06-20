package gorm_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testMock() (*gorm.DB, sqlmock.Sqlmock, *domain.Account) {
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

	accountMock, _ := domain.NewAccount(20.75)

	return gormDB, mock, accountMock
}
