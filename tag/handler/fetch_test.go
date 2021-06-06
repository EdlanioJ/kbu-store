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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FetchTagsHandler(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		fetchTagUsecase := new(mocks.FetchTagsUsecase)
		sort := "total"
		page := 1
		limit := 5
		fetchTagUsecase.
			On("Exec", mock.Anything, sort, page, limit).
			Return(nil, int64(0), errors.New("Unexpexted Error")).
			Once()

		fetchTags := handler.NewFetchTags(fetchTagUsecase)

		app.Get("/", fetchTags.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/?sort=%s&page=%d&limit=%d", sort, page, limit), nil)

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
		fetchTagUsecase := new(mocks.FetchTagsUsecase)
		sort := "total"
		page := 1
		limit := 5
		fetchTagUsecase.
			On("Exec", mock.Anything, sort, page, limit).
			Return(tags, int64(1), nil).
			Once()

		fetchTags := handler.NewFetchTags(fetchTagUsecase)

		app.Get("/", fetchTags.Handler)
		req := httptest.NewRequest("GET", fmt.Sprintf("/?sort=%s&page=%d&limit=%d", sort, page, limit), nil)

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.Header.Get("X-total"), "1")
		is.Equal(res.StatusCode, fiber.StatusOK)
	})
}
