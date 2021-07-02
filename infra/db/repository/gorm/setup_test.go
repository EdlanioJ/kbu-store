package gorm_test

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/domain"
	uuid "github.com/satori/go.uuid"
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

func getStore() *domain.Store {
	c := new(domain.Category)
	c.ID = uuid.NewV4().String()

	store := &domain.Store{
		Base: domain.Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
		},
		Name:        "Store 001",
		Description: "Store description 001",
		Status:      domain.StoreStatusActive,
		ExternalID:  uuid.NewV4().String(),
		AccountID:   uuid.NewV4().String(),
		Category:    c,
		Tags:        []string{"tag 001", "tag 002"},
		Position: domain.Position{
			Lat: -8.8867698,
			Lng: 13.4771186,
		},
	}
	return store
}
