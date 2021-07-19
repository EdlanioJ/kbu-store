package usecase_test

import (
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	uuid "github.com/satori/go.uuid"
)

func getStore() *domain.Store {
	c := new(domain.Category)
	c.ID = uuid.NewV4().String()
	store := &domain.Store{
		Base: domain.Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "Store 001",
		Description: "store description 001",
		Status:      domain.StoreStatusPending,
		UserID:      uuid.NewV4().String(),
		AccountID:   uuid.NewV4().String(),
		Category:    c,
	}

	return store
}

func getCategory() *domain.Category {

	mockStorType, _ := domain.NewCategory("Store type 001")

	return mockStorType
}
