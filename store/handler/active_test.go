package handler_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/store/handler"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ActiveHandler(t *testing.T) {
	t.Run("fail on id validation", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		id := "invalid_id"
		activeHandler := handler.NewActiveHandler(nil)

		app.Options("/:id/active", activeHandler.Handler)
		req := httptest.NewRequest(fiber.MethodOptions, fmt.Sprintf("/%s/active", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		activeUsecase := new(mocks.ActivateStoreUsecase)
		id := uuid.NewV4().String()

		activeUsecase.On("Exec", mock.Anything, id).Return(errors.New("failed"))
		activeHandler := handler.NewActiveHandler(activeUsecase)

		app.Options("/:id/active", activeHandler.Handler)
		req := httptest.NewRequest(fiber.MethodOptions, fmt.Sprintf("/%s/active", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)

		activeUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		activeUsecase := new(mocks.ActivateStoreUsecase)
		id := uuid.NewV4().String()

		activeUsecase.On("Exec", mock.Anything, id).Return(nil)
		activeHandler := handler.NewActiveHandler(activeUsecase)

		app.Options("/:id/active", activeHandler.Handler)
		req := httptest.NewRequest(fiber.MethodOptions, fmt.Sprintf("/%s/active", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusNoContent)

		activeUsecase.AssertExpectations(t)
	})
}
