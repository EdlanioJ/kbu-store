package handler_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/tag/handler"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FetchTagsByCategoryHandler(t *testing.T) {
	t.Run("fail on category id validation", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		categoryID := "invalid_id"

		sort := "total"
		page := 1
		limit := 5

		fetchTagsByCategory := handler.NewFetchTagsByCategory(nil)

		app.Get("/category/:category", fetchTagsByCategory.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/category/%s?sort=%s&page=%d&limit=%d", categoryID, sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		categoryID := uuid.NewV4().String()
		sort := "total"
		page := 1
		limit := 5

		fetchTagByCategoryUsecase := new(mocks.FetchTagsByCategoryUsecase)
		fetchTagByCategoryUsecase.
			On("Exec", mock.Anything, categoryID, sort, limit, page).
			Return(nil, int64(0), errors.New("Unexpexted Error")).
			Once()

		fetchTagsByCategory := handler.NewFetchTagsByCategory(fetchTagByCategoryUsecase)

		app.Get("/category/:category", fetchTagsByCategory.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/category/%s?sort=%s&page=%d&limit=%d", categoryID, sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "0")
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, tag := mockTest()

		tags := make([]*domain.Tag, 0)
		tags = append(tags, tag)
		categoryID := uuid.NewV4().String()
		sort := "total"
		page := 1
		limit := 5

		fetchTagByCategoryUsecase := new(mocks.FetchTagsByCategoryUsecase)
		fetchTagByCategoryUsecase.
			On("Exec", mock.Anything, categoryID, sort, limit, page).
			Return(tags, int64(1), nil).
			Once()

		fetchTagsByCategory := handler.NewFetchTagsByCategory(fetchTagByCategoryUsecase)

		app.Get("/category/:category", fetchTagsByCategory.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/category/%s?sort=%s&page=%d&limit=%d", categoryID, sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "1")
		is.Equal(res.StatusCode, fiber.StatusOK)
	})
}
