package handler_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EdlanioJ/kbu-store/category/handler"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateHandler(t *testing.T) {
	t.Run("fail on parser body", func(t *testing.T) {
		is := assert.New(t)
		app, _ := testMock()

		createHandler := handler.NewCreateHandler(nil)
		app.Post("/", createHandler.Handler)
		req := httptest.NewRequest(fiber.MethodPost, "/", strings.NewReader(`{error: this is wrong}`))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusUnprocessableEntity)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := testMock()
		createUsecase := new(mocks.CreateCategoryUsecase)

		cr := new(dto.CreateCategoryRequest)
		cr.Name = "Store type 001"

		c, err := json.Marshal(cr)
		is.NoError(err)

		createUsecase.On("Add", mock.Anything, cr.Name).Return(errors.New("failed"))
		createHandler := handler.NewCreateHandler(createUsecase)
		app.Post("/", createHandler.Handler)
		req := httptest.NewRequest(fiber.MethodPost, "/", strings.NewReader(string(c)))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, _ := testMock()
		createUsecase := new(mocks.CreateCategoryUsecase)

		cr := new(dto.CreateCategoryRequest)
		cr.Name = "Store type 001"

		c, err := json.Marshal(cr)
		is.NoError(err)

		createUsecase.On("Add", mock.Anything, cr.Name).Return(nil)
		createHandler := handler.NewCreateHandler(createUsecase)
		app.Post("/", createHandler.Handler)
		req := httptest.NewRequest(fiber.MethodPost, "/", strings.NewReader(string(c)))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusCreated)
	})
}
