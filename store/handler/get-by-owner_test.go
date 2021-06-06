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

func Test_GetByOwnerHandler(t *testing.T) {
	t.Run("fail on id validation", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		id := "invalid_id"
		ownerID := uuid.NewV4().String()

		getByOwnerHandler := handler.NewGetByOwnerHandler(nil)

		app.Get("/:id/:owner", getByOwnerHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", id, ownerID), nil)

		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on owner validation", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		id := uuid.NewV4().String()
		ownerID := "invalid_uuid"

		getByOwnerHandler := handler.NewGetByOwnerHandler(nil)

		app.Get("/:id/:owner", getByOwnerHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", id, ownerID), nil)

		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		getByOwnerUsecase := new(mocks.GetStoreByOwnerUsecase)

		id := uuid.NewV4().String()
		ownerID := uuid.NewV4().String()

		getByOwnerUsecase.On("Exec", mock.Anything, id, ownerID).Return(nil, errors.New("Unexpexted Error"))
		getByOwnerHandler := handler.NewGetByOwnerHandler(getByOwnerUsecase)

		app.Get("/:id/:owner", getByOwnerHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", id, ownerID), nil)
		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)

		getByOwnerUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, store := mockTest()
		getByOwnerUsecase := new(mocks.GetStoreByOwnerUsecase)

		id := uuid.NewV4().String()
		ownerID := uuid.NewV4().String()

		getByOwnerUsecase.On("Exec", mock.Anything, id, ownerID).Return(store, nil)
		getByOwnerHandler := handler.NewGetByOwnerHandler(getByOwnerUsecase)

		app.Get("/:id/:owner", getByOwnerHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", id, ownerID), nil)
		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusOK)
		getByOwnerUsecase.AssertExpectations(t)
	})
}
