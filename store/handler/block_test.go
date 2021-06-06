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

func Test_BlockHandler(t *testing.T) {
	t.Run("fail on id validation", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		id := "invalid_id"
		blockHandler := handler.NewBlockHandler(nil)

		app.Get("/:id/block", blockHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/block", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		blockUsecase := new(mocks.BlockStoreUsecase)
		id := uuid.NewV4().String()

		blockUsecase.On("Exec", mock.Anything, id).Return(domain.ErrBlocked)
		blockHandler := handler.NewBlockHandler(blockUsecase)

		app.Get("/:id/block", blockHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/block", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)

		blockUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		blockUsecase := new(mocks.BlockStoreUsecase)
		id := uuid.NewV4().String()

		blockUsecase.On("Exec", mock.Anything, id).Return(nil)
		blockHandler := handler.NewBlockHandler(blockUsecase)

		app.Get("/:id/block", blockHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s/block", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusNoContent)

		blockUsecase.AssertExpectations(t)
	})
}
