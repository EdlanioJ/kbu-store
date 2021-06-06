package handler_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/store/handler"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetByIDHandler(t *testing.T) {
	t.Run("fail on validation id", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		id := "invalid_id"

		getByIDHandler := handler.NewGetByIDHandler(nil)

		app.Get("/:id", getByIDHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s", id), nil)

		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		getByIDUsecase := new(mocks.GetStoreByIDUsecase)
		id := uuid.NewV4().String()

		getByIDUsecase.On("Exec", mock.Anything, id).Return(nil, domain.ErrNotFound)
		getByIDHandler := handler.NewGetByIDHandler(getByIDUsecase)

		app.Get("/:id", getByIDHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s", id), nil)

		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusNotFound)
		getByIDUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, store := mockTest()

		getByIDUsecase := new(mocks.GetStoreByIDUsecase)
		id := uuid.NewV4().String()

		getByIDUsecase.On("Exec", mock.Anything, id).Return(store, nil)
		getByIDHandler := handler.NewGetByIDHandler(getByIDUsecase)

		app.Get("/:id", getByIDHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s", id), nil)

		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusOK)
		getByIDUsecase.AssertExpectations(t)
	})
}
