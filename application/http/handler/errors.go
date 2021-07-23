package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Message string `json:"msg"`
}

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, e error) error {
		return c.Status(getStatusCode(e)).JSON(ErrorResponse{
			Message: e.Error(),
		})
	}
}

func getStatusCode(err error) int {
	logrus.Error(err)

	if _, ok := err.(govalidator.Errors); ok {
		return fiber.StatusBadRequest
	}

	switch err {
	case domain.ErrNotFound,
		gorm.ErrRecordNotFound:
		return fiber.StatusNotFound
	case domain.ErrBadParam:
		return fiber.StatusBadRequest
	case domain.ErrActived,
		domain.ErrBlocked,
		domain.ErrInactived,
		domain.ErrIsPending:
		return fiber.StatusConflict
	}
	return fiber.StatusInternalServerError
}
