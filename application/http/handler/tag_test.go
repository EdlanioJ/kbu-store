package handler_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/EdlanioJ/kbu-store/application/http/handler"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_TagHandler_GetAll(t *testing.T) {
	type args struct {
		sort  string
		page  int
		limit int
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(tagUsecase *mocks.TagUsecase)
		checkResponse func(t *testing.T, err error, total string, statusCode int)
	}{
		{
			name: "fail",
			args: args{
				sort:  "created_at",
				page:  1,
				limit: 5,
			},
			builtSts: func(tagUsecase *mocks.TagUsecase) {
				tagUsecase.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, total string, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
				assert.Equal(t, total, "0")
			},
		},

		{
			name: "success",
			args: args{
				sort:  "created_at",
				page:  1,
				limit: 5,
			},
			builtSts: func(tagUsecase *mocks.TagUsecase) {
				tag := &domain.Tag{
					Name:  "tag001",
					Count: 10,
				}
				tags := make([]*domain.Tag, 0)
				tags = append(tags, tag)
				tagUsecase.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(tags, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, total string, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
				assert.Equal(t, total, "1")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tagUsecase := new(mocks.TagUsecase)
			tc.builtSts(tagUsecase)
			app := fiber.New()
			handler.NewTagRoutes(app, tagUsecase)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/tags/?sort=%s&page=%d&limit=%d", tc.args.sort, tc.args.page, tc.args.limit), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			total := res.Header.Get("X-total")
			tc.checkResponse(t, err, total, res.StatusCode)
			tagUsecase.AssertExpectations(t)
		})
	}
}

func Test_TagHandler_GetAllByCategory(t *testing.T) {
	type args struct {
		categoryId string
		sort       string
		page       int
		limit      int
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(tagUsecase *mocks.TagUsecase)
		checkResponse func(t *testing.T, err error, total string, statusCode int)
	}{
		{
			name: "fail on category id validation",
			args: args{
				categoryId: "invald_id",
				sort:       "created_at",
				page:       1,
				limit:      5,
			},
			builtSts: func(_ *mocks.TagUsecase) {},
			checkResponse: func(t *testing.T, err error, _ string, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			args: args{
				categoryId: uuid.NewV4().String(),
				sort:       "created_at",
				page:       1,
				limit:      5,
			},
			builtSts: func(tagUsecase *mocks.TagUsecase) {
				tagUsecase.On("GetAllByCategory", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, total string, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
				assert.Equal(t, total, "0")
			},
		},
		{
			name: "success",
			args: args{
				categoryId: uuid.NewV4().String(),
				sort:       "created_at",
				page:       1,
				limit:      5,
			},
			builtSts: func(tagUsecase *mocks.TagUsecase) {
				tag := &domain.Tag{
					Name:  "tag001",
					Count: 10,
				}
				tags := make([]*domain.Tag, 0)
				tags = append(tags, tag)
				tagUsecase.On("GetAllByCategory", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(tags, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, total string, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
				assert.Equal(t, total, "1")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tagUsecase := new(mocks.TagUsecase)
			tc.builtSts(tagUsecase)
			app := fiber.New()
			handler.NewTagRoutes(app, tagUsecase)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/tags/category/%s?sort=%s&page=%d&limit=%d", tc.args.categoryId, tc.args.sort, tc.args.page, tc.args.limit), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			total := res.Header.Get("X-total")
			tc.checkResponse(t, err, total, res.StatusCode)
			tagUsecase.AssertExpectations(t)
		})
	}
}
