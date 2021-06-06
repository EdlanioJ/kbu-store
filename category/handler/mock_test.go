package handler_test

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/gofiber/fiber/v2"
)

func testMock() (*fiber.App, *domain.Category) {
	storType, _ := domain.NewCategory("Store type 001")
	app := fiber.New()

	return app, storType
}
