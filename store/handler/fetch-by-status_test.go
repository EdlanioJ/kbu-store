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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FetchByStatusHandler(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		fetchByStatusUsecase := new(mocks.FetchStoreByStatusUsecase)

		sort := "created_at"
		page := 1
		limit := 5
		status := domain.StoreStatusPending

		fetchByStatusUsecase.On("Exec", mock.Anything, status, sort, limit, page).Return(nil, int64(0), errors.New("Unexpexted Error"))
		fetchByStatusHandler := handler.NewFetchByStatusHandler(fetchByStatusUsecase)
		app.Get("/:status", fetchByStatusHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s?sort=%s&page=%d&limit=%d", status, sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "0")
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)

		fetchByStatusUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, store := mockTest()
		fetchByStatusUsecase := new(mocks.FetchStoreByStatusUsecase)

		stores := make([]*domain.Store, 0)
		stores = append(stores, store)

		sort := "created_at"
		page := 1
		limit := 5
		status := domain.StoreStatusPending

		fetchByStatusUsecase.On("Exec", mock.Anything, status, sort, limit, page).Return(stores, int64(1), nil).Once()
		fetchByStatusHandler := handler.NewFetchByStatusHandler(fetchByStatusUsecase)
		app.Get("/:status", fetchByStatusHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s?sort=%s&page=%d&limit=%d", status, sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "1")
		is.Equal(res.StatusCode, fiber.StatusOK)

		fetchByStatusUsecase.AssertExpectations(t)
	})
}
