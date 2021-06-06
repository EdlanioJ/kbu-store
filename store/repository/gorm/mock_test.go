package gorm_test

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu-store/domain"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testMock() (*gorm.DB, sqlmock.Sqlmock, *domain.Store) {
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

	a := new(domain.Account)
	st := new(domain.Category)
	a.ID = uuid.NewV4().String()
	st.ID = uuid.NewV4().String()

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
		Account:     a,
		Category:    st,
		Tags:        []string{"tag 001", "tag 002"},
		Position: domain.Position{
			Lat: -8.8867698,
			Lng: 13.4771186,
		},
	}

	return gormDB, mock, store
}
