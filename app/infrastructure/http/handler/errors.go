package handler

import (
	"errors"
	"strings"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Message string      `json:"message"`
	Field   string      `json:"field,omitempty"`
	Value   interface{} `json:"value,omitempty" swaggertype:"object"`
}

type HttpError struct {
	Status int
	Error  interface{}
}

func errorHandler(c *fiber.Ctx, err error) error {
	httpError := getHttpError(err)
	return c.Status(httpError.Status).JSON(httpError.Error)
}

func getHttpError(err error) HttpError {
	switch {
	case errors.Is(err, domain.ErrNotFound),
		errors.Is(err, gorm.ErrRecordNotFound):
		return HttpError{
			Status: fiber.StatusNotFound,
			Error:  ErrorResponse{Message: domain.ErrNotFound.Error()},
		}
	case errors.Is(err, domain.ErrActived),
		errors.Is(err, domain.ErrBlocked),
		errors.Is(err, domain.ErrInactived),
		errors.Is(err, domain.ErrIsPending):
		return HttpError{
			Status: fiber.StatusConflict,
			Error:  ErrorResponse{Message: err.Error()},
		}
	case strings.Contains(strings.ToLower(err.Error()), "json"):
		return HttpError{
			Status: fiber.StatusBadRequest,
			Error:  ErrorResponse{Message: domain.ErrBadRequest.Error()},
		}
	default:
		if err, ok := err.(validator.ValidationErrors); ok {
			var errors []ErrorResponse
			for _, e := range err {
				errors = append(errors, ErrorResponse{
					Message: e.Error(),
					Field:   e.Field(),
					Value:   e.Value(),
				})
			}
			return HttpError{
				Status: fiber.StatusBadRequest,
				Error:  errors,
			}
		}

		return HttpError{
			Status: fiber.StatusInternalServerError,
			Error:  ErrorResponse{Message: domain.ErrInternal.Error()},
		}
	}
}
