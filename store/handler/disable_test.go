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

func Test_DisableHandle(t *testing.T) {
	t.Run("fail on id validation", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		id := "invalid_id"
		disableHandler := handler.NewDisableHandler(nil)

		app.Get("/:id/disable", disableHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/disable", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		disableUsecase := new(mocks.DisableStoreUsecase)
		id := uuid.NewV4().String()

		disableUsecase.On("Exec", mock.Anything, id).Return(errors.New("failed"))
		disableHandler := handler.NewDisableHandler(disableUsecase)

		app.Get("/:id/disable", disableHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/disable", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)

		disableUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		disableUsecase := new(mocks.DisableStoreUsecase)
		id := uuid.NewV4().String()

		disableUsecase.On("Exec", mock.Anything, id).Return(nil)
		disableHandler := handler.NewDisableHandler(disableUsecase)

		app.Get("/:id/disable", disableHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/disable", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusNoContent)

		disableUsecase.AssertExpectations(t)
	})
}
