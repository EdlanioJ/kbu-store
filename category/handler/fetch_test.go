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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FetchHandler(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)
		app, _ := testMock()
		fetchUsecase := new(mocks.FetchCategoryUsecase)

		sort := "created_at"
		page := 1
		limit := 5

		fetchUsecase.On("Exec", mock.Anything, sort, limit, page).Return(nil, int64(0), errors.New("Unexpexted Error"))
		fetchHandler := handler.NewFetchHandler(fetchUsecase)

		app.Get("/", fetchHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/?sort=%s&page=%d&limit=%d", sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "0")
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, Category := testMock()
		fetchUsecase := new(mocks.FetchCategoryUsecase)

		categories := make([]*domain.Category, 0)
		categories = append(categories, Category)
		sort := "created_at"
		page := 1
		limit := 5

		fetchUsecase.On("Exec", mock.Anything, sort, limit, page).Return(categories, int64(1), nil)
		fetchHandler := handler.NewFetchHandler(fetchUsecase)

		app.Get("/", fetchHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/?sort=%s&page=%d&limit=%d", sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "1")
		is.Equal(res.StatusCode, fiber.StatusOK)
	})
}
