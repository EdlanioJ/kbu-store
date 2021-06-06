package handler_test

import (
	"errors"
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

func Test_FetchByTypeHandler(t *testing.T) {
	t.Run("fail on validation category", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		sort := "created_at"
		page := 1
		limit := 5
		categoryID := "invald_uuid"

		fetchByTypeHandler := handler.NewFetchByTypeHandler(nil)

		app.Get("/:category", fetchByTypeHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s?sort=%s&page=%d&limit=%d", categoryID, sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		sort := "created_at"
		page := 1
		limit := 5
		categoryID := uuid.NewV4().String()

		fetchByTypeUsecase := new(mocks.FetchStoreByTypeUsecase)

		fetchByTypeUsecase.On("Exec", mock.Anything, categoryID, sort, limit, page).Return(nil, int64(0), errors.New("Unexpexted Error"))
		fetchByTypeHandler := handler.NewFetchByTypeHandler(fetchByTypeUsecase)

		app.Get("/:category", fetchByTypeHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s?sort=%s&page=%d&limit=%d", categoryID, sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "0")
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)

		fetchByTypeUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, store := mockTest()

		sort := "created_at"
		page := 1
		limit := 5
		categoryID := uuid.NewV4().String()

		stores := make([]*domain.Store, 0)
		stores = append(stores, store)
		fetchByTypeUsecase := new(mocks.FetchStoreByTypeUsecase)

		fetchByTypeUsecase.On("Exec", mock.Anything, categoryID, sort, limit, page).Return(stores, int64(1), nil)
		fetchByTypeHandler := handler.NewFetchByTypeHandler(fetchByTypeUsecase)

		app.Get("/:category", fetchByTypeHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s?sort=%s&page=%d&limit=%d", categoryID, sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "1")
		is.Equal(res.StatusCode, fiber.StatusOK)

		fetchByTypeUsecase.AssertExpectations(t)
	})
}
