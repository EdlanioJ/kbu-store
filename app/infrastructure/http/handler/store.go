package handler

import (
	"fmt"
	"strconv"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/validators"
	"github.com/gofiber/fiber/v2"
)

type storeHandler struct {
	storeUsecase domain.StoreUsecase
}

func NewStoreHandler(us domain.StoreUsecase) *storeHandler {
	return &storeHandler{
		storeUsecase: us,
	}
}

// @Summary Create  a new store
// @Description Create store
// @Tags stores
// @Accept json
// @Produce json
// @Param category body CreateStoreRequest true "Create store"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /stores [post]
func (h *storeHandler) Store(c *fiber.Ctx) error {
	ctx := c.Context()
	cr := new(CreateStoreRequest)
	if err := c.BodyParser(cr); err != nil {
		return err
	}
	err := h.storeUsecase.Store(ctx, cr.Name, cr.Description, cr.CategoryID, cr.UserID, cr.Tags, cr.Lat, cr.Lng)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusCreated)
}

// @Summary Get all stores
// @Description Get a list of stores
// @Tags stores
// @Accept json
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param sort query string false "Sort" default(created_at DESC)
// @Success 200 {array} domain.Store
// @Failure 500 {object} ErrorResponse
// @Router /stores [get]
func (h *storeHandler) Index(c *fiber.Ctx) error {
	sort := c.Query("sort")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	ctx := c.Context()

	list, total, err := h.storeUsecase.Index(ctx, sort, limit, page)
	c.Response().Header.Add("X-total", fmt.Sprint(total))
	if err != nil {
		return err
	}
	return c.JSON(list)
}

// @Summary Get stores by id
// @Description Get one stores by id
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 200 {object} domain.Store
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id} [get]
func (h *storeHandler) Get(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return err
	}

	res, err := h.storeUsecase.Get(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// @Summary Activate stores
// @Description Activate one stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id}/activate [patch]
func (h *storeHandler) Activate(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return err
	}

	err = h.storeUsecase.Active(ctx, id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Block stores
// @Description Block one stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id}/block [patch]
func (h *storeHandler) Block(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return err
	}

	err = h.storeUsecase.Block(ctx, id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Disable stores
// @Description Disable one stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id}/disable [patch]
func (h *storeHandler) Disable(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return err
	}

	err = h.storeUsecase.Disable(ctx, id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Delete stores
// @Description Delete one stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id} [delete]
func (h *storeHandler) Delete(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return err
	}

	err = h.storeUsecase.Delete(ctx, id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary update store
// @Description Activate one stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Param category body UpdateStoreRequest true "Create store"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/active/{id} [patch]
func (h *storeHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()
	reqBody := new(UpdateStoreRequest)
	id := c.Params("id")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return err
	}

	if err := c.BodyParser(reqBody); err != nil {
		return err
	}

	store := reqBody.ToDomainStore()
	store.ID = id

	err = h.storeUsecase.Update(ctx, store)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
