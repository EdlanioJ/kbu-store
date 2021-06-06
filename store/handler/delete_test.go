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

func Test_DeleteHandler(t *testing.T) {
	t.Run("fail on id validation", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		id := "invalid_id"
		deleteHandler := handler.NewDeleteHandler(nil)

		app.Delete("/:id", deleteHandler.Handler)
		req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%s", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		id := uuid.NewV4().String()
		deleteUsecase := new(mocks.DeleteStoreUsecase)

		deleteUsecase.On("Exec", mock.Anything, id).Return(errors.New("failed")).Once()
		deleteHandler := handler.NewDeleteHandler(deleteUsecase)

		app.Delete("/:id", deleteHandler.Handler)
		req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%s", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		id := uuid.NewV4().String()
		deleteUsecase := new(mocks.DeleteStoreUsecase)

		deleteUsecase.On("Exec", mock.Anything, id).Return(nil).Once()
		deleteHandler := handler.NewDeleteHandler(deleteUsecase)

		app.Delete("/:id", deleteHandler.Handler)
		req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%s", id), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusNoContent)
	})
}
