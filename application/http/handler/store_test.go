package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/EdlanioJ/kbu-store/application/http/handler"
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getStore() *domain.Store {
	mockStorType, _ := domain.NewCategory("Store type 001")

	storeMock := &domain.Store{
		Base: domain.Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "Store 001",
		Description: "store description 001",
		Status:      domain.StoreStatusPending,
		ExternalID:  uuid.NewV4().String(),
		AccountID:   uuid.NewV4().String(),
		Category:    mockStorType,
	}

	return storeMock
}
func Test_StoreHandler_Create(t *testing.T) {
	cr := handler.CreateStoreRequest{
		Name:        "Store 001",
		Description: "Store description",
		CategoryID:  uuid.NewV4().String(),
		ExternalID:  uuid.NewV4().String(),
		Tags:        []string{"tag001", "tag002"},
		Lat:         -8.8867698,
		Lng:         13.4771186,
	}
	c, _ := json.Marshal(cr)
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "fail on parser body",
			arg:      "{error: this is wrong}",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusUnprocessableEntity)
			},
		},
		{
			name: "fail on usecase",
			arg:  string(c),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Create", mock.Anything, cr.Name, cr.Description, cr.CategoryID, cr.ExternalID, cr.Tags, cr.Lat, cr.Lng).Return(errors.New("failed")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "success",
			arg:  string(c),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Create", mock.Anything, cr.Name, cr.Description, cr.CategoryID, cr.ExternalID, cr.Tags, cr.Lat, cr.Lng).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusCreated)
			},
		},
	}

	for _, tc := range testCases {
		storeUsecase := new(mocks.StoreUsecase)
		tc.builtSts(storeUsecase)
		app := fiber.New()
		handler.NewStoreRoute(app, storeUsecase)
		req := httptest.NewRequest(fiber.MethodPost, "/stores/", strings.NewReader(tc.arg))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		tc.checkResponse(t, err, res.StatusCode)
		storeUsecase.AssertExpectations(t)
	}
}

func Test_StoreHandler_GetAll(t *testing.T) {
	type args struct {
		sort  string
		page  int
		limit int
	}

	mockArgs := args{
		sort:  "created_at",
		page:  1,
		limit: 5,
	}

	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int, total string)
	}{
		{
			name: "fail",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("GetAll", mock.Anything, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
				assert.Equal(t, total, "0")
			},
		},
		{
			name: "success",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				store := getStore()
				stores := make([]*domain.Store, 0)
				stores = append(stores, store)
				storeUsecase.On("GetAll", mock.Anything, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(stores, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
				assert.Equal(t, total, "1")
			},
		},
	}

	for _, tc := range testCases {
		storeUsecase := new(mocks.StoreUsecase)
		tc.builtSts(storeUsecase)
		app := fiber.New()
		handler.NewStoreRoute(app, storeUsecase)
		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/stores/?sort=%s&page=%d&limit=%d", tc.args.sort, tc.args.page, tc.args.limit), nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		total := res.Header.Get("X-total")
		tc.checkResponse(t, err, res.StatusCode, total)
		storeUsecase.AssertExpectations(t)
	}
}

func Test_StoreHandler_GetAllByCategory(t *testing.T) {
	type args struct {
		categoryId string
		sort       string
		page       int
		limit      int
	}

	mockArgs := args{
		categoryId: uuid.NewV4().String(),
		sort:       "created_at",
		page:       1,
		limit:      5,
	}

	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int, total string)
	}{
		{
			name: "fail on category id validation",
			args: args{
				categoryId: "invalid_id",
				sort:       "created_at",
				page:       1,
				limit:      5,
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int, _ string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("GetAllByCategory", mock.Anything, mockArgs.categoryId, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
				assert.Equal(t, total, "0")
			},
		},
		{
			name: "success",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				store := getStore()
				stores := make([]*domain.Store, 0)
				stores = append(stores, store)
				storeUsecase.On("GetAllByCategory", mock.Anything, mockArgs.categoryId, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(stores, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
				assert.Equal(t, total, "1")
			},
		},
	}

	for _, tc := range testCases {
		storeUsecase := new(mocks.StoreUsecase)
		tc.builtSts(storeUsecase)
		app := fiber.New()
		handler.NewStoreRoute(app, storeUsecase)
		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/stores/category/%s?sort=%s&page=%d&limit=%d", tc.args.categoryId, tc.args.sort, tc.args.page, tc.args.limit), nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		total := res.Header.Get("X-total")
		tc.checkResponse(t, err, res.StatusCode, total)
		storeUsecase.AssertExpectations(t)
	}
}

func Test_StoreHandler_GetAllBCloseLocation(t *testing.T) {
	type args struct {
		lat      float64
		lng      float64
		sort     string
		page     int
		limit    int
		distance int
		status   string
	}

	mockArgs := args{
		lat:      -8.8867698,
		lng:      13.4771186,
		sort:     "created_at",
		page:     1,
		limit:    5,
		distance: 10,
		status:   domain.StoreStatusPending,
	}

	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int, total string)
	}{
		{
			name: "fail on latitude validation",
			args: args{
				lat:    -300.8867698,
				lng:    13.4771186,
				sort:   "created_at",
				page:   1,
				limit:  5,
				status: domain.StoreStatusActive,
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int, _ string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on longitude validation",
			args: args{
				lat:    -8.8867698,
				lng:    333.4771186,
				sort:   "created_at",
				page:   1,
				limit:  5,
				status: domain.StoreStatusBlock,
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int, _ string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("GetAllByByCloseLocation", mock.Anything, mockArgs.lat, mockArgs.lng, mockArgs.distance, mockArgs.status, mockArgs.limit, mockArgs.page, mockArgs.sort).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
				assert.Equal(t, total, "0")
			},
		},
		{
			name: "success",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				store := getStore()
				stores := make([]*domain.Store, 0)
				stores = append(stores, store)
				storeUsecase.On("GetAllByByCloseLocation", mock.Anything, mockArgs.lat, mockArgs.lng, mockArgs.distance, mockArgs.status, mockArgs.limit, mockArgs.page, mockArgs.sort).Return(stores, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
				assert.Equal(t, total, "1")
			},
		},
	}

	for _, tc := range testCases {
		storeUsecase := new(mocks.StoreUsecase)
		tc.builtSts(storeUsecase)
		app := fiber.New()
		handler.NewStoreRoute(app, storeUsecase)
		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/stores/location/%s/status/%s?sort=%s&page=%d&limit=%d&distance=%d", fmt.Sprintf("@%v,%v", tc.args.lat, tc.args.lng), tc.args.status, tc.args.sort, tc.args.page, tc.args.limit, tc.args.distance), nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		total := res.Header.Get("X-total")
		tc.checkResponse(t, err, res.StatusCode, total)
		storeUsecase.AssertExpectations(t)
	}
}

func Test_StoreHandler_GetAllByOwner(t *testing.T) {
	type args struct {
		ownerId string
		sort    string
		page    int
		limit   int
	}

	mockArgs := args{
		ownerId: uuid.NewV4().String(),
		sort:    "created_at",
		page:    1,
		limit:   5,
	}

	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int, total string)
	}{
		{
			name: "fail on owner id validation",
			args: args{
				ownerId: "invalid_id",
				sort:    "created_at",
				page:    1,
				limit:   5,
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int, _ string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("GetAllByOwner", mock.Anything, mockArgs.ownerId, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
				assert.Equal(t, total, "0")
			},
		},
		{
			name: "success",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				store := getStore()
				stores := make([]*domain.Store, 0)
				stores = append(stores, store)
				storeUsecase.On("GetAllByOwner", mock.Anything, mockArgs.ownerId, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(stores, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
				assert.Equal(t, total, "1")
			},
		},
	}

	for _, tc := range testCases {
		storeUsecase := new(mocks.StoreUsecase)
		tc.builtSts(storeUsecase)
		app := fiber.New()
		handler.NewStoreRoute(app, storeUsecase)
		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/stores/owner/%s?sort=%s&page=%d&limit=%d", tc.args.ownerId, tc.args.sort, tc.args.page, tc.args.limit), nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		total := res.Header.Get("X-total")
		tc.checkResponse(t, err, res.StatusCode, total)
		storeUsecase.AssertExpectations(t)
	}
}

func Test_StoreHandler_GetAllByStatus(t *testing.T) {
	type args struct {
		status string
		sort   string
		page   int
		limit  int
	}

	mockArgs := args{
		status: domain.StoreStatusActive,
		sort:   "created_at",
		page:   1,
		limit:  5,
	}

	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int, total string)
	}{
		{
			name: "fail",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("GetAllByStatus", mock.Anything, mockArgs.status, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
				assert.Equal(t, total, "0")
			},
		},
		{
			name: "success",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				store := getStore()
				stores := make([]*domain.Store, 0)
				stores = append(stores, store)
				storeUsecase.On("GetAllByStatus", mock.Anything, mockArgs.status, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(stores, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
				assert.Equal(t, total, "1")
			},
		},
	}

	for _, tc := range testCases {
		storeUsecase := new(mocks.StoreUsecase)
		tc.builtSts(storeUsecase)
		app := fiber.New()
		handler.NewStoreRoute(app, storeUsecase)
		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/stores/status/%s?sort=%s&page=%d&limit=%d", tc.args.status, tc.args.sort, tc.args.page, tc.args.limit), nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		total := res.Header.Get("X-total")
		tc.checkResponse(t, err, res.StatusCode, total)
		storeUsecase.AssertExpectations(t)
	}
}

func Test_StoreHandler_GetAllByTags(t *testing.T) {
	type args struct {
		tags  string
		sort  string
		page  int
		limit int
	}

	mockArgs := args{
		tags:  "tag001,tag002",
		sort:  "created_at",
		page:  1,
		limit: 5,
	}

	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int, total string)
	}{
		{
			name: "fail",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("GetAllByTags", mock.Anything, mock.Anything, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
				assert.Equal(t, total, "0")
			},
		},
		{
			name: "success",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				store := getStore()
				stores := make([]*domain.Store, 0)
				stores = append(stores, store)
				storeUsecase.On("GetAllByTags", mock.Anything, mock.Anything, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(stores, int64(1), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
				assert.Equal(t, total, "1")
			},
		},
	}

	for _, tc := range testCases {
		storeUsecase := new(mocks.StoreUsecase)
		tc.builtSts(storeUsecase)
		app := fiber.New()
		handler.NewStoreRoute(app, storeUsecase)
		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/stores/tags/?tags=%s&sort=%s&page=%d&limit=%d", tc.args.tags, tc.args.sort, tc.args.page, tc.args.limit), nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		total := res.Header.Get("X-total")
		tc.checkResponse(t, err, res.StatusCode, total)
		storeUsecase.AssertExpectations(t)
	}
}

func Test_StoreHandler_GetById(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "fail on id validation",
			arg:      "invalid_id",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, domain.ErrNotFound).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNotFound)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(getStore(), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeUsecase := new(mocks.StoreUsecase)
			tc.builtSts(storeUsecase)
			app := fiber.New()
			handler.NewStoreRoute(app, storeUsecase)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/stores/%s", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			tc.checkResponse(t, err, res.StatusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_GetByIdAndOwner(t *testing.T) {
	type args struct {
		id    string
		owner string
	}
	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name: "fail on id validation",
			args: args{
				id:    "invalid_id",
				owner: uuid.NewV4().String(),
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on owner id validation",
			args: args{
				id:    uuid.NewV4().String(),
				owner: "invalid_id",
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			args: args{
				id:    uuid.NewV4().String(),
				owner: uuid.NewV4().String(),
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("GetByIdAndOwner", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, domain.ErrNotFound).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNotFound)
			},
		},
		{
			name: "success",
			args: args{
				id:    uuid.NewV4().String(),
				owner: uuid.NewV4().String(),
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("GetByIdAndOwner", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(getStore(), nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusOK)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeUsecase := new(mocks.StoreUsecase)
			tc.builtSts(storeUsecase)
			app := fiber.New()
			handler.NewStoreRoute(app, storeUsecase)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/stores/%s/owner/%s", tc.args.id, tc.args.owner), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			tc.checkResponse(t, err, res.StatusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Active(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "fail on id validation",
			arg:      "invalid_id",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Active", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrActived).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Active", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNoContent)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeUsecase := new(mocks.StoreUsecase)
			tc.builtSts(storeUsecase)
			app := fiber.New()
			handler.NewStoreRoute(app, storeUsecase)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/stores/%s/activate", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			tc.checkResponse(t, err, res.StatusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Block(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "fail on id validation",
			arg:      "invalid_id",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Block", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrBlocked).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Block", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNoContent)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeUsecase := new(mocks.StoreUsecase)
			tc.builtSts(storeUsecase)
			app := fiber.New()
			handler.NewStoreRoute(app, storeUsecase)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/stores/%s/block", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			tc.checkResponse(t, err, res.StatusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Disable(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "fail on id validation",
			arg:      "invalid_id",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Disable", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrBlocked).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Disable", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNoContent)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeUsecase := new(mocks.StoreUsecase)
			tc.builtSts(storeUsecase)
			app := fiber.New()
			handler.NewStoreRoute(app, storeUsecase)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/stores/%s/disable", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			tc.checkResponse(t, err, res.StatusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Delete(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "fail on id validation",
			arg:      "invalid_id",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on usecase",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("Unextpected Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "success",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNoContent)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeUsecase := new(mocks.StoreUsecase)
			tc.builtSts(storeUsecase)
			app := fiber.New()
			handler.NewStoreRoute(app, storeUsecase)
			req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/stores/%s", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			tc.checkResponse(t, err, res.StatusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Update(t *testing.T) {
	type args struct {
		id      string
		request string
	}

	ur := handler.UpdateStoreRequest{
		Name:        "store 002",
		Description: "description 002",
		CategoryID:  uuid.NewV4().String(),
		Tags:        []string{"tag002", "tag003"},
		Lat:         -8.8867698,
		Lng:         13.4771186,
	}

	c, _ := json.Marshal(ur)
	testCases := []struct {
		name          string
		args          args
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name: "fail on id validation",
			args: args{
				id:      "invalid_id",
				request: "",
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusBadRequest)
			},
		},
		{
			name: "fail on body parser",
			args: args{
				id:      uuid.NewV4().String(),
				request: "{error: this is wrong}",
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusUnprocessableEntity)
			},
		},
		{
			name: "fail on usecase",
			args: args{
				id:      uuid.NewV4().String(),
				request: string(c),
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Update", mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "fail on usecase",
			args: args{
				id:      uuid.NewV4().String(),
				request: string(c),
			},
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusNoContent)
			},
		},
	}

	for _, tc := range testCases {
		storeUsecase := new(mocks.StoreUsecase)
		tc.builtSts(storeUsecase)
		app := fiber.New()
		handler.NewStoreRoute(app, storeUsecase)
		req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/stores/%s", tc.args.id), strings.NewReader(tc.args.request))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		tc.checkResponse(t, err, res.StatusCode)
		storeUsecase.AssertExpectations(t)
	}
}
