package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
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

func Test_UpdateHandler(t *testing.T) {
	t.Run("fail on id validation", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		id := "invalid_id"
		updateHandler := handler.NewUpdateHandler(nil)
		app.Put("/:id", updateHandler.Handler)
		req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/%s", id), nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusBadRequest)
	})

	t.Run("fail on body parser", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		id := uuid.NewV4().String()

		updateHandler := handler.NewUpdateHandler(nil)
		app.Put("/:id", updateHandler.Handler)
		req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/%s", id), strings.NewReader(`{error: this is wrong}`))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusUnprocessableEntity)
	})

	t.Run("fail on usecase", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		id := uuid.NewV4().String()

		ur := dto.UpdateStoreRequest{
			Name:        "store 002",
			Description: "description 002",
			CategoryID:  uuid.NewV4().String(),
			Tags:        []string{"tag002", "tag003"},
			Lat:         -8.8867698,
			Lng:         13.4771186,
		}

		c, err := json.Marshal(ur)
		is.NoError(err)

		updateUsecase := new(mocks.UpdateStoreUsecase)

		updateUsecase.On("Exec", mock.Anything, mock.Anything).Return(errors.New("failed"))
		updateHandler := handler.NewUpdateHandler(updateUsecase)
		app.Put("/:id", updateHandler.Handler)
		req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/%s", id), strings.NewReader(string(c)))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusInternalServerError)
		updateUsecase.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		is := assert.New(t)
		app, _ := mockTest()
		id := uuid.NewV4().String()

		ur := dto.UpdateStoreRequest{
			Name:        "store 002",
			Description: "description 002",
			CategoryID:  uuid.NewV4().String(),
			Tags:        []string{"tag002", "tag003"},
			Lat:         -8.8867698,
			Lng:         13.4771186,
		}

		c, err := json.Marshal(ur)
		is.NoError(err)

		updateUsecase := new(mocks.UpdateStoreUsecase)

		updateUsecase.On("Exec", mock.Anything, mock.Anything).Return(nil)
		updateHandler := handler.NewUpdateHandler(updateUsecase)
		app.Put("/:id", updateHandler.Handler)
		req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/%s", id), strings.NewReader(string(c)))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)

		is.NoError(err)
		is.Equal(res.StatusCode, fiber.StatusNoContent)
		updateUsecase.AssertExpectations(t)
	})
}
