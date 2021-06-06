package handler_test

import (
	"time"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
)

func mockTest() (*fiber.App, *domain.Store) {
	mockStorType, _ := domain.NewCategory("Store type 001")

	account := new(domain.Account)

	account.ID = uuid.NewV4().String()

	app := fiber.New()
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
		Account:     account,
		Category:    mockStorType,
	}

	return app, storeMock
}
