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

func Test_FetchByLocationAndStatusHandler(t *testing.T) {
	t.Run("fail on invalid latitude", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		fetchByLocation := handler.NewFetchByLocationAndStatusHandler(nil)

		sort := "created_at"
		page := 1
		limit := 5

		location := "@-300.8867698,13.4771186"
		distance := 10
		status := domain.StoreStatusPending

		app.Get("/location/:location/status/:status", fetchByLocation.Handler)

		req := httptest.NewRequest("GET", fmt.Sprintf("/location/%s/status/%s?sort=%s&page=%d&limit=%d&distance=%d", location, status, sort, page, limit, distance), nil)
		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on invalid longuitude", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		fetchByLocation := handler.NewFetchByLocationAndStatusHandler(nil)

		sort := "created_at"
		page := 1
		limit := 5

		location := "@-8.8867698,333.4771186"
		distance := 10
		status := domain.StoreStatusPending

		app.Get("/location/:location/status/:status", fetchByLocation.Handler)

		req := httptest.NewRequest("GET", fmt.Sprintf("/location/%s/status/%s?sort=%s&page=%d&limit=%d&distance=%d", location, status, sort, page, limit, distance), nil)
		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		fetchByLocationUsecase := new(mocks.FetchStoreByCloseLocationUsecase)
		fetchByLocation := handler.NewFetchByLocationAndStatusHandler(fetchByLocationUsecase)

		sort := "created_at"
		page := 1
		limit := 5

		location := "@-8.8867698,13.4771186"
		lat := -8.8867698
		lng := 13.4771186
		distance := 10
		status := domain.StoreStatusPending

		fetchByLocationUsecase.
			On("Exec", mock.Anything, lat, lng, distance, status, limit, page, sort).
			Return(nil, int64(0), errors.New("Unexpexted Error"))
		app.Get("/location/:location/status/:status", fetchByLocation.Handler)

		req := httptest.NewRequest("GET", fmt.Sprintf("/location/%s/status/%s?sort=%s&page=%d&limit=%d&distance=%d", location, status, sort, page, limit, distance), nil)
		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "0")
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)
		fetchByLocationUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, store := mockTest()

		fetchByLocationUsecase := new(mocks.FetchStoreByCloseLocationUsecase)
		fetchByLocation := handler.NewFetchByLocationAndStatusHandler(fetchByLocationUsecase)

		sort := "created_at"
		page := 1
		limit := 5

		location := "@-8.8867698,13.4771186"
		lat := -8.8867698
		lng := 13.4771186
		distance := 10
		status := domain.StoreStatusPending
		stores := make([]*domain.Store, 0)
		stores = append(stores, store)
		fetchByLocationUsecase.
			On("Exec", mock.Anything, lat, lng, distance, status, limit, page, sort).
			Return(stores, int64(1), nil)
		app.Get("/location/:location/status/:status", fetchByLocation.Handler)

		req := httptest.NewRequest("GET", fmt.Sprintf("/location/%s/status/%s?sort=%s&page=%d&limit=%d&distance=%d", location, status, sort, page, limit, distance), nil)
		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "1")
		is.Equal(res.StatusCode, fiber.StatusOK)
		fetchByLocationUsecase.AssertExpectations(t)
	})
}
