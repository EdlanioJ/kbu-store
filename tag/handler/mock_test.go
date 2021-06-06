package handler_test

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/gofiber/fiber/v2"
)

func mockTest() (*fiber.App, *domain.Tag) {
	app := fiber.New()
	return app, &domain.Tag{
		Name:  "tag001",
		Count: 4,
	}
}
