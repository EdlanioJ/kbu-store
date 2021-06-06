package handler_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/EdlanioJ/kbu-store/category/handler"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetByStatusHandler(t *testing.T) {
	t.Run("fail on id validation", func(t *testing.T) {
		is := assert.New(t)
		app, _ := testMock()

		id := "invalid_id"
		status := domain.CategoryStatusPending

		getByStatusHandler := handler.NewGetByStatusHandler(nil)

		app.Get("/:id/:status", getByStatusHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", id, status), nil)

		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := testMock()
		getByStatusUsecase := new(mocks.GetCategoryByStautsRepository)

		id := uuid.NewV4().String()
		status := domain.CategoryStatusPending

		getByStatusUsecase.On("Exec", mock.Anything, id, status).Return(nil, errors.New("Unexpexted Error"))
		getByStatusHandler := handler.NewGetByStatusHandler(getByStatusUsecase)

		app.Get("/:id/:status", getByStatusHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", id, status), nil)
		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)

		getByStatusUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, store := testMock()
		getByStatusUsecase := new(mocks.GetCategoryByStautsRepository)
		id := uuid.NewV4().String()
		status := domain.CategoryStatusPending

		getByStatusUsecase.On("Exec", mock.Anything, id, status).Return(store, nil)
		getByStatusHandler := handler.NewGetByStatusHandler(getByStatusUsecase)

		app.Get("/:id/:status", getByStatusHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", id, status), nil)
		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusOK)
		getByStatusUsecase.AssertExpectations(t)
	})
}
