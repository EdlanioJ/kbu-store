package handler_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EdlanioJ/kbu-store/domain/mocks"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/EdlanioJ/kbu-store/store/handler"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateHandler(t *testing.T) {
	t.Run("fail on parser body", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()

		createHandler := handler.NewCreateHandler(nil)
		app.Post("/", createHandler.Handler)
		req := httptest.NewRequest(fiber.MethodPost, "/", strings.NewReader(`{error: this is wrong}`))

		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusUnprocessableEntity)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		createUsecase := new(mocks.CreateStoreUsecase)

		cr := dto.CreateStoreRequest{
			Name:        "Store 001",
			Description: "Store description",
			CategoryID:  uuid.NewV4().String(),
			ExternalID:  uuid.NewV4().String(),
			Tags:        []string{"tag001", "tag002"},
			Lat:         -8.8867698,
			Lng:         13.4771186,
		}

		c, err := json.Marshal(cr)
		is.NoError(err)

		createUsecase.On("Add", mock.Anything, cr.Name, cr.Description, cr.CategoryID, cr.ExternalID, cr.Tags, cr.Lat, cr.Lng).Return(errors.New("failed")).Once()
		createHandler := handler.NewCreateHandler(createUsecase)
		app.Post("/", createHandler.Handler)
		req := httptest.NewRequest(fiber.MethodPost, "/", strings.NewReader(string(c)))

		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		createUsecase := new(mocks.CreateStoreUsecase)

		cr := dto.CreateStoreRequest{
			Name:        "Store 001",
			Description: "Store description",
			CategoryID:  uuid.NewV4().String(),
			ExternalID:  uuid.NewV4().String(),
			Tags:        []string{"tag001", "tag002"},
			Lat:         -8.8867698,
			Lng:         13.4771186,
		}

		c, err := json.Marshal(cr)
		is.NoError(err)

		createUsecase.On("Add", mock.Anything, cr.Name, cr.Description, cr.CategoryID, cr.ExternalID, cr.Tags, cr.Lat, cr.Lng).Return(nil).Once()
		createHandler := handler.NewCreateHandler(createUsecase)
		app.Post("/", createHandler.Handler)
		req := httptest.NewRequest(fiber.MethodPost, "/", strings.NewReader(string(c)))

		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusCreated)
	})
}
