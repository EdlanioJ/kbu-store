package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/http/handler"
	"github.com/EdlanioJ/kbu-store/app/utils/mocks"
	"github.com/EdlanioJ/kbu-store/app/utils/sample"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_StoreHandler_Store(t *testing.T) {
	cr := sample.NewCreateStoreRequest()
	invalidCreateStore := sample.NewCreateStoreRequest()
	invalidCreateStore.Name = ""
	c, err := json.Marshal(cr)
	assert.NoError(t, err)
	ic, err := json.Marshal(invalidCreateStore)
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		arg           string
		statusCode    int
		prepare       func(storeUsecase *mocks.StoreUsecase)
		checkResponse func(t *testing.T, err error, statusCode int)
	}{
		{
			name:       "failure_parser_body",
			arg:        `{error: this is wrong}`,
			statusCode: fiber.StatusBadRequest,
		},
		{
			name:       "failure_validate_create_request",
			arg:        string(ic),
			statusCode: fiber.StatusBadRequest,
		},
		{
			name:       "failure_usecase_returns_error",
			arg:        string(c),
			statusCode: fiber.StatusInternalServerError,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Store", mock.Anything, cr).Return((domain.ErrInternal)).Once()
			},
		},
		{
			name:       "success",
			arg:        string(c),
			statusCode: fiber.StatusCreated,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Store", mock.Anything, cr).Return(nil).Once()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeUsecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(storeUsecase)
			}
			app := fiber.New()
			validator := validator.New()
			handler := handler.NewStoreHandler(storeUsecase, validator)
			app.Post("/", handler.Store)
			req := httptest.NewRequest(fiber.MethodPost, "/", strings.NewReader(tc.arg))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, tc.statusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Index(t *testing.T) {
	args := sample.NewHttpListReq()
	testCases := []struct {
		name       string
		args       sample.HttpListRequest
		statusCode int
		prepare    func(storeUsecase *mocks.StoreUsecase)
	}{
		{
			name:       "failure",
			args:       args,
			statusCode: fiber.StatusInternalServerError,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Index", mock.Anything, args.Sort, args.Limit, args.Page).Return(nil, int64(0), errors.New("Unexpected Error")).Once()
			},
		},
		{
			name:       "success",
			args:       args,
			statusCode: fiber.StatusOK,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				store := sample.NewStore()
				stores := make(domain.Stores, 0)
				stores = append(stores, store)
				storeUsecase.On("Index", mock.Anything, args.Sort, args.Limit, args.Page).Return(stores, int64(1), nil).Once()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeUsecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(storeUsecase)
			}
			app := fiber.New()
			validator := validator.New()
			handler := handler.NewStoreHandler(storeUsecase, validator)
			app.Get("/", handler.Index)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/?sort=%s&page=%d&limit=%d", tc.args.Sort, tc.args.Page, tc.args.Limit), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, tc.statusCode)
			if res.StatusCode == fiber.StatusOK {
				total := res.Header.Get("X-total")
				assert.Equal(t, total, "1")
			}
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Get(t *testing.T) {
	testCases := []struct {
		name       string
		arg        string
		statusCode int
		prepare    func(storeUsecase *mocks.StoreUsecase)
	}{
		{
			name:       "failure_invalid_id",
			arg:        "invalid_id",
			statusCode: fiber.StatusBadRequest,
		},
		{
			name:       "failure_usecase_returns_error",
			arg:        uuid.NewV4().String(),
			statusCode: fiber.StatusNotFound,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, domain.ErrNotFound).Once()
			},
		},
		{
			name:       "should succeed",
			arg:        uuid.NewV4().String(),
			statusCode: fiber.StatusOK,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(sample.NewStore(), nil).Once()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeUsecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(storeUsecase)
			}
			app := fiber.New()
			validator := validator.New()
			handler := handler.NewStoreHandler(storeUsecase, validator)
			app.Get("/:id", handler.Get)
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%s", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, tc.statusCode)
		})
	}
}

func Test_StoreHandler_Activate(t *testing.T) {
	testCases := []struct {
		name       string
		arg        string
		statusCode int
		prepare    func(storeUsecase *mocks.StoreUsecase)
	}{
		{
			name:       "failure_invalid_id",
			arg:        "invalid_id",
			statusCode: fiber.StatusBadRequest,
		},
		{
			name:       "failure_usercase_returns_error",
			arg:        uuid.NewV4().String(),
			statusCode: fiber.StatusConflict,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Active", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrActived).Once()
			},
		},
		{
			name:       "success",
			arg:        uuid.NewV4().String(),
			statusCode: fiber.StatusNoContent,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Active", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeUsecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(storeUsecase)
			}
			app := fiber.New()
			validator := validator.New()
			handler := handler.NewStoreHandler(storeUsecase, validator)
			app.Patch("/:id/activate", handler.Activate)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/%s/activate", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, tc.statusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Block(t *testing.T) {
	testCases := []struct {
		name       string
		arg        string
		statusCode int
		prepare    func(storeUsecase *mocks.StoreUsecase)
	}{
		{
			name:       "failure_usecase_returns_error",
			arg:        "invalid_id",
			statusCode: fiber.StatusBadRequest,
		},
		{
			name:       "should fail if usecase returns an error",
			arg:        uuid.NewV4().String(),
			statusCode: fiber.StatusConflict,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Block", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrBlocked).Once()
			},
		},
		{
			name:       "success",
			arg:        uuid.NewV4().String(),
			statusCode: fiber.StatusNoContent,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Block", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeUsecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(storeUsecase)
			}
			app := fiber.New()
			validator := validator.New()
			handler := handler.NewStoreHandler(storeUsecase, validator)
			app.Patch("/:id/block", handler.Block)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/%s/block", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, tc.statusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Disable(t *testing.T) {
	testCases := []struct {
		name       string
		arg        string
		statusCode int
		prepare    func(storeUsecase *mocks.StoreUsecase)
	}{
		{
			name:       "failure_invalid_id",
			arg:        "invalid_id",
			statusCode: fiber.StatusBadRequest,
		},
		{
			name:       "failure_usecase_returns_error",
			arg:        uuid.NewV4().String(),
			statusCode: fiber.StatusConflict,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Disable", mock.Anything, mock.AnythingOfType("string")).Return(domain.ErrBlocked).Once()
			},
		},
		{
			name:       "success",
			arg:        uuid.NewV4().String(),
			statusCode: fiber.StatusNoContent,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Disable", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeUsecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(storeUsecase)
			}
			app := fiber.New()
			validator := validator.New()
			handler := handler.NewStoreHandler(storeUsecase, validator)
			app.Patch("/:id/disable", handler.Disable)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/%s/disable", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, tc.statusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Delete(t *testing.T) {
	testCases := []struct {
		name       string
		arg        string
		statusCode int
		prepare    func(storeUsecase *mocks.StoreUsecase)
	}{
		{
			name:       "failure_invalid_id",
			arg:        "invalid_id",
			statusCode: fiber.StatusBadRequest,
		},
		{
			name:       "failure_usecase_returns_error",
			arg:        uuid.NewV4().String(),
			statusCode: fiber.StatusInternalServerError,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("Unextpected Error")).Once()
			},
		},
		{
			name:       "success",
			arg:        uuid.NewV4().String(),
			statusCode: fiber.StatusNoContent,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeUsecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(storeUsecase)
			}
			app := fiber.New()
			validator := validator.New()
			handler := handler.NewStoreHandler(storeUsecase, validator)
			app.Delete("/:id", handler.Delete)
			req := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%s", tc.arg), nil)
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, tc.statusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}

func Test_StoreHandler_Update(t *testing.T) {
	ur := sample.NewUpdateStoreRequest()

	c, _ := json.Marshal(ur)
	testCases := []struct {
		name       string
		id         string
		request    string
		statusCode int
		prepare    func(storeUsecase *mocks.StoreUsecase)
	}{
		{
			name:       "failure_invalid_od",
			id:         "invalid_id",
			request:    "",
			statusCode: fiber.StatusBadRequest,
		},
		{
			name:       "failure_parser_body",
			id:         uuid.NewV4().String(),
			request:    "{error: this is wrong}",
			statusCode: fiber.StatusBadRequest,
		},
		{
			name:       "failure_validate_returns_error",
			id:         "invalid_id",
			request:    string(c),
			statusCode: fiber.StatusBadRequest,
		},
		{
			name:       "failure_usecase_returns_error",
			id:         uuid.NewV4().String(),
			request:    string(c),
			statusCode: fiber.StatusInternalServerError,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Update", mock.Anything, mock.Anything).Return(errors.New("Unextpected Error"))
			},
		},
		{
			name:       "success",
			id:         uuid.NewV4().String(),
			request:    string(c),
			statusCode: fiber.StatusNoContent,
			prepare: func(storeUsecase *mocks.StoreUsecase) {
				storeUsecase.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			storeUsecase := new(mocks.StoreUsecase)
			if tc.prepare != nil {
				tc.prepare(storeUsecase)
			}
			app := fiber.New()
			validator := validator.New()
			handler := handler.NewStoreHandler(storeUsecase, validator)
			app.Patch("/:id", handler.Update)
			req := httptest.NewRequest(fiber.MethodPatch, fmt.Sprintf("/%s", tc.id), strings.NewReader(tc.request))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, tc.statusCode)
			storeUsecase.AssertExpectations(t)
		})
	}
}
