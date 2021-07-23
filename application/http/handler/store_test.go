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
	store := &domain.Store{
		Base: domain.Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "Store 001",
		Description: "store description 001",
		Status:      domain.StoreStatusPending,
		UserID:      uuid.NewV4().String(),
		AccountID:   uuid.NewV4().String(),
		CategoryID:  uuid.NewV4().String(),
	}

	return store
}
func Test_StoreHandler_Store(t *testing.T) {
	cr := handler.CreateStoreRequest{
		Name:        "Store 001",
		Description: "Store description",
		CategoryID:  uuid.NewV4().String(),
		UserID:      uuid.NewV4().String(),
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
			name:     "should fail if parser body returns an err",
			arg:      "{error: this is wrong}",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should fail if usecase returns an error",
			arg:  string(c),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Store", mock.Anything, cr.Name, cr.Description, cr.CategoryID, cr.UserID, cr.Tags, cr.Lat, cr.Lng).Return(errors.New("failed")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should succeed",
			arg:  string(c),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Store", mock.Anything, cr.Name, cr.Description, cr.CategoryID, cr.UserID, cr.Tags, cr.Lat, cr.Lng).Return(nil).Once()
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
		handler := handler.NewStoreHandler(storeUsecase)
		app.Post("/", handler.Store)
		req := httptest.NewRequest(fiber.MethodPost, "/", strings.NewReader(tc.arg))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		tc.checkResponse(t, err, res.StatusCode)
		storeUsecase.AssertExpectations(t)
	}
}

func Test_StoreHandler_Index(t *testing.T) {
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
			name: "should fail if usecase returns an error",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Index", mock.Anything, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(nil, int64(0), errors.New("Unexpexted Error")).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int, total string) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
				assert.Equal(t, total, "0")
			},
		},
		{
			name: "should succeed",
			args: mockArgs,
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				store := getStore()
				stores := make(domain.Stores, 0)
				stores = append(stores, store)
				storeUsecase.On("Index", mock.Anything, mockArgs.sort, mockArgs.limit, mockArgs.page).Return(stores, int64(1), nil).Once()
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
		handler := handler.NewStoreHandler(storeUsecase)
		app.Get("/", handler.Index)
		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/?sort=%s&page=%d&limit=%d", tc.args.sort, tc.args.page, tc.args.limit), nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		total := res.Header.Get("X-total")
		tc.checkResponse(t, err, res.StatusCode, total)
		storeUsecase.AssertExpectations(t)
	}
}

func Test_StoreHandler_Get(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "should fail if id validation returns an error",
			arg:      "invalid_id",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should fail if usecase returns an error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, domain.ErrNotFound).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should succeed",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(getStore(), nil).Once()
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
			handler := handler.NewStoreHandler(storeUsecase)
			app.Get("/:id", handler.Get)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%s", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			tc.checkResponse(t, err, res.StatusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Activate(t *testing.T) {
	testCases := []struct {
		name          string
		arg           string
		builtSts      func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:     "should fail if id validation returns an error",
			arg:      "invalid_id",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should fail if usecase returns an error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Active", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrActived).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should succeed",
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
			handler := handler.NewStoreHandler(storeUsecase)
			app.Patch("/:id/activate", handler.Activate)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/%s/activate", tc.arg), nil)
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
			name:     "should fail if id validation returns an error",
			arg:      "invalid_id",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should fail if usecase returns an error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Block", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrBlocked).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should succeed",
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
			handler := handler.NewStoreHandler(storeUsecase)
			app.Patch("/:id/block", handler.Block)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/%s/block", tc.arg), nil)
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
			name:     "should fail if id validation returns an error",
			arg:      "invalid_id",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should fail if usecase returns an error",
			arg:  uuid.NewV4().String(),
			builtSts: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Disable", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrBlocked).Once()
			},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should succeed",
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
			handler := handler.NewStoreHandler(storeUsecase)
			app.Patch("/:id/disable", handler.Disable)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/%s/disable", tc.arg), nil)
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
			name:     "should fail if id validation returns an error",
			arg:      "invalid_id",
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should fail if usecase returns an error",
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
			name: "should succeed",
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
			handler := handler.NewStoreHandler(storeUsecase)
			app.Delete("/:id", handler.Delete)
			req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%s", tc.arg), nil)
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
			name: "should fail if id validation returns an error",
			args: args{
				id:      "invalid_id",
				request: "",
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should fail if body parser returns an error",
			args: args{
				id:      uuid.NewV4().String(),
				request: "{error: this is wrong}",
			},
			builtSts: func(_ *mocks.StoreUsecase) {},
			checkResponse: func(t *testing.T, err error, statusCode int) {
				assert.NoError(t, err)
				assert.Equal(t, statusCode, fiber.StatusInternalServerError)
			},
		},
		{
			name: "should fail if usecase returns an error",
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
		handler := handler.NewStoreHandler(storeUsecase)
		app.Patch("/:id", handler.Update)
		req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/%s", tc.args.id), strings.NewReader(tc.args.request))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		tc.checkResponse(t, err, res.StatusCode)
		storeUsecase.AssertExpectations(t)
	}
}
