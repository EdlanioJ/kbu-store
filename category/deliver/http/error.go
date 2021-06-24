package http

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Message string `json:"msg"`
}

func getStatusCode(err error) int {
	logrus.Error(err)

	switch err {
	case domain.ErrNotFound:
		return fiber.StatusNotFound
	case domain.ErrActived,
		domain.ErrBlocked,
		domain.ErrBadParam,
		domain.ErrActived,
		domain.ErrInactived,
		domain.ErrIsPending:
		return fiber.StatusBadRequest
	}
	return fiber.StatusInternalServerError
}
