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

func Test_FetchByTagsHandler(t *testing.T) {

	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		fetchByTagsUsecase := new(mocks.FetchStoreByTagsUsecase)

		sort := "created_at"
		page := 1
		limit := 5
		tags := "tag001,tag002"
		fetchByTagsUsecase.On("Exec", mock.Anything, mock.Anything, sort, limit, page).Return(nil, int64(0), errors.New("Unexpexted Error"))
		fetchByTagsHandler := handler.NewFetchByTagsHandler(fetchByTagsUsecase)

		app.Get("/", fetchByTagsHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/?tags=%s&sort=%s&page=%d&limit=%d", tags, sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "0")
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)

		fetchByTagsUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, store := mockTest()
		fetchByTagsUsecase := new(mocks.FetchStoreByTagsUsecase)

		sort := "created_at"
		page := 1
		limit := 5
		tags := "tag001,tag002"
		stores := make([]*domain.Store, 0)
		stores = append(stores, store)

		fetchByTagsUsecase.On("Exec", mock.Anything, mock.Anything, sort, limit, page).Return(stores, int64(1), nil)
		fetchByTagsHandler := handler.NewFetchByTagsHandler(fetchByTagsUsecase)

		app.Get("/", fetchByTagsHandler.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/?tags=%s&sort=%s&page=%d&limit=%d", tags, sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "1")
		is.Equal(res.StatusCode, fiber.StatusOK)

		fetchByTagsUsecase.AssertExpectations(t)
	})
}
