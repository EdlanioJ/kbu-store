package usecase_test

import (
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	uuid "github.com/satori/go.uuid"
)

type entityMock struct {
	store    *domain.Store
	storType *domain.Category
}

func testMock() *entityMock {
	storType, _ := domain.NewCategory("Store type 001")
	account := new(domain.Account)
	account.ID = uuid.NewV4().String()

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
		Account:     account,
		Category:    storType,
	}

	return &entityMock{
		store:    store,
		storType: storType,
	}
}
