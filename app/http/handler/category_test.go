package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EdlanioJ/kbu-store/app/http/handler"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getCategory() *domain.Category {
	category, _ := domain.NewCategory("Store type 001")

	return category
}

func Test_CategoryHandler_Create(t *testing.T) {
	cr := new(handler.CreateCategoryRequest)
	cr.Name = "Store type 001"
	c, _ := json.Marshal(cr)

	testCases := []struct {
		name          string
		arg           string
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "fail on parser body",
			arg:      "{error: this is wrong}",
			builtSts: func(_ *mocks.CategoryUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusUnprocessableEntity)
			},
		},
		{
			name: "fail on usecase",
			arg:  string(c),
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Create", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("failed")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "success",
			arg:  string(c),
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Create", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusCreated)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryUsecase := new(mocks.CategoryUsecase)

			tc.builtSts(categoryUsecase)
			app := fiber.New()
			handler.NewCategoryRoutes(app, categoryUsecase)
			req := httptest.NewRequest(fiber.MethodPost, "/categories/", strings.NewReader(tc.arg))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			tc.checkResponse(t, err, res.StatusCode)
			categoryUsecase.AssertExpectations(t)
		})
	}
}

func Test_CategoryHandler_GetById(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "fail on validation id",
			arg:      "invald_id",
			builtSts: func(_ *mocks.CategoryUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, domain.ErrNotFound).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNotFound)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(getCategory(), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryUsecase := new(mocks.CategoryUsecase)
			tc.builtSts(categoryUsecase)
			app := fiber.New()
			handler.NewCategoryRoutes(app, categoryUsecase)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/categories/%s", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			tc.checkResponse(t, err, res.StatusCode)
			categoryUsecase.AssertExpectations(t)
		})
	}
}

func Test_CategoryHandler_GetByIdAndStatus(t *testing.T) {
	type args struct {
		id     string
		status string
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name: "fail on id validation",
			args: args{
				id:     "invalid_id",
				status: domain.CategoryStatusActive,
			},
			builtSts: func(_ *mocks.CategoryUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			args: args{
				id:     uuid.NewV4().String(),
				status: domain.CategoryStatusActive,
			},
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("GetByIdAndStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "success",
			args: args{
				id:     uuid.NewV4().String(),
				status: domain.CategoryStatusActive,
			},
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("GetByIdAndStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(getCategory(), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryUsecase := new(mocks.CategoryUsecase)
			tc.builtSts(categoryUsecase)
			app := fiber.New()
			handler.NewCategoryRoutes(app, categoryUsecase)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/categories/%s/status/%s", tc.args.id, tc.args.status), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			tc.checkResponse(t, err, res.StatusCode)
			categoryUsecase.AssertExpectations(t)
		})
	}
}
func Test_CategoryHandler_GetAll(t *testing.T) {
	type args struct {
		sort  string
		page  int
		limit int
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, err error, total string, statusCode int)
	}{
		{
			name: "fail",
			args: args{
				sort:  "created_at",
				page:  1,
				limit: 5,
			},
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
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
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categories := make([]*domain.Category, 0)
				categories = append(categories, getCategory())
				categoryUsecase.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(categories, int64(1), nil).Once()
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
			categoryUsecase := new(mocks.CategoryUsecase)
			tc.builtSts(categoryUsecase)
			app := fiber.New()
			handler.NewCategoryRoutes(app, categoryUsecase)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/categories/?sort=%s&page=%d&limit=%d", tc.args.sort, tc.args.page, tc.args.limit), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			total := res.Header.Get("X-total")
			tc.checkResponse(t, err, total, res.StatusCode)
			categoryUsecase.AssertExpectations(t)
		})
	}
}

func Test_CategoryHandler_GetAllByStatus(t *testing.T) {
	type args struct {
		status string
		sort   string
		page   int
		limit  int
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, err error, total string, statusCode int)
	}{
		{
			name: "fail",
			args: args{
				status: domain.CategoryStatusActive,
				sort:   "created_at",
				page:   1,
				limit:  5,
			},
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("GetAllByStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, int64(0), domain.ErrBadParam).Once()
			},
			checkResponse: func(t *testing.T, err error, total string, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
				assert.Equal(t, total, "0")
			},
		},

		{
			name: "success",
			args: args{
				status: domain.CategoryStatusInactive,
				sort:   "created_at",
				page:   1,
				limit:  5,
			},
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categories := make([]*domain.Category, 0)
				categories = append(categories, getCategory())
				categoryUsecase.On("GetAllByStatus", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(categories, int64(1), nil).Once()
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
			categoryUsecase := new(mocks.CategoryUsecase)
			tc.builtSts(categoryUsecase)
			app := fiber.New()
			handler.NewCategoryRoutes(app, categoryUsecase)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/categories/status/%s?sort=%s&page=%d&limit=%d", tc.args.status, tc.args.sort, tc.args.page, tc.args.limit), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			total := res.Header.Get("X-total")
			tc.checkResponse(t, err, total, res.StatusCode)
			categoryUsecase.AssertExpectations(t)
		})
	}
}

func Test_CategoryHandler_Activate(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "fail on validation id",
			arg:      "invald_id",
			builtSts: func(_ *mocks.CategoryUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Activate", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrNotFound).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNotFound)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Activate", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNoContent)
			},
		},
	}

	for _, tc := range testCases {
		categoryUsecase := new(mocks.CategoryUsecase)
		tc.builtSts(categoryUsecase)
		app := fiber.New()
		handler.NewCategoryRoutes(app, categoryUsecase)
		req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/categories/%s/activate", tc.arg), nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		tc.checkResponse(t, err, res.StatusCode)
		categoryUsecase.AssertExpectations(t)
	}
}

func Test_CategoryHandler_Disable(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(categoryUsecase *mocks.CategoryUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "fail on validation id",
			arg:      "invald_id",
			builtSts: func(_ *mocks.CategoryUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Disable", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrNotFound).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNotFound)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(categoryUsecase *mocks.CategoryUsecase) {
				categoryUsecase.On("Disable", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNoContent)
			},
		},
	}

	for _, tc := range testCases {
		categoryUsecase := new(mocks.CategoryUsecase)
		tc.builtSts(categoryUsecase)
		app := fiber.New()
		handler.NewCategoryRoutes(app, categoryUsecase)
		req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/categories/%s/disable", tc.arg), nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		tc.checkResponse(t, err, res.StatusCode)
		categoryUsecase.AssertExpectations(t)
	}
}
