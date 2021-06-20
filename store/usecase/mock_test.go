package usecase_test

import (
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	uuid "github.com/satori/go.uuid"
)

type entityMock struct {
	store    *domain.Store
	category *domain.Category
}

func testMock() *entityMock {
	category, _ := domain.NewCategory("Store type 001")

	store := &domain.Store{
		Base: domain.Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "Store 001",
		Description: "store description 001",
		Status:      domain.StoreStatusPending,
		ExternalID:  uuid.NewV4().String(),
		AccountID:   uuid.NewV4().String(),
		Category:    category,
	}

	return &entityMock{
		store:    store,
		category: category,
	}
}
