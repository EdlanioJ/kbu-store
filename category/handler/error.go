package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func getStatusCode(err error) int {
	logrus.Error(err)

	switch err {
	case domain.ErrNotFound:
		return fiber.StatusNotFound
	case domain.ErrBadParam:
		return fiber.StatusBadRequest
	}
	return fiber.StatusInternalServerError
}
