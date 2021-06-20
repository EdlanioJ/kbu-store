package http_test

import (
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	uuid "github.com/satori/go.uuid"
)

func mockTest() *domain.Store {
	mockStorType, _ := domain.NewCategory("Store type 001")

	storeMock := &domain.Store{
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
		Category:    mockStorType,
	}

	return storeMock
}
