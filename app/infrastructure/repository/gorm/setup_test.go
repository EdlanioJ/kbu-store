package gorm_test

import (
	"github.com/DATA-DOG/go-sqlmock"
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
