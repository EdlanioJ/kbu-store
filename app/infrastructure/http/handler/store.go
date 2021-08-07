package handler

import (
	"fmt"
	"strconv"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

type storeHandler struct {
	storeUsecase domain.StoreUsecase
	validate     *validator.Validate
}

func NewStoreHandler(usecase domain.StoreUsecase, validate *validator.Validate) *storeHandler {
	return &storeHandler{
		storeUsecase: usecase,
		validate:     validate,
	}
}

// @Summary Create store
// @Description Create new store
// @Tags stores
// @Accept json
// @Produce json
// @Param category body domain.CreateStoreRequest true "Create store"
// @Success 201 {string} string "Created"
// @Failure 400 {array} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /stores [post]
func (h *storeHandler) Store(c *fiber.Ctx) error {
	span, ctx := opentracing.StartSpanFromContext(c.Context(), "StoreHandler.Store")
	defer span.Finish()
	createRequests.Inc()

	cr := new(domain.CreateStoreRequest)
	if err := c.BodyParser(cr); err != nil {
		log.
			WithContext(ctx).
			Errorf("c.BodyParser: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	if err := h.validate.StructCtx(ctx, cr); err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.StructCtx: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	err := h.storeUsecase.Store(ctx, cr)
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Store: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	successRequests.Inc()
	return c.SendStatus(fiber.StatusCreated)

}

// @Summary Index store
// @Description Get list of stores
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
	span, ctx := opentracing.StartSpanFromContext(c.Context(), "StoreHandler.Index")
	defer span.Finish()
	indexRequests.Inc()

	sort := c.Query("sort")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	list, total, err := h.storeUsecase.Index(ctx, sort, limit, page)
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Index: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	c.Response().Header.Add("X-total", fmt.Sprint(total))
	successRequests.Inc()
	return c.JSON(list)
}

// @Summary Get stores
// @Description Get a stores by id
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 200 {object} domain.Store
// @Failure 500 {object} ErrorResponse
// @Failure 400 {array} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id} [get]
func (h *storeHandler) Get(c *fiber.Ctx) error {
	span, ctx := opentracing.StartSpanFromContext(c.Context(), "StoreHandler.Get")
	defer span.Finish()
	getRequests.Inc()

	id := c.Params("id")

	err := h.validate.VarCtx(ctx, id, "uuid4")
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.VarCtx: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	res, err := h.storeUsecase.Get(ctx, id)
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Get: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	successRequests.Inc()
	return c.JSON(res)
}

// @Summary Activate stores
// @Description Activate a stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {array} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id}/activate [patch]
func (h *storeHandler) Activate(c *fiber.Ctx) error {
	span, ctx := opentracing.StartSpanFromContext(c.Context(), "StoreHandler.Activate")
	defer span.Finish()
	ativateRequests.Inc()

	id := c.Params("id")

	err := h.validate.VarCtx(ctx, id, "uuid4")
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.VarCtx: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	err = h.storeUsecase.Active(ctx, id)
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Active: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	successRequests.Inc()
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Block stores
// @Description Block a stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {array} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id}/block [patch]
func (h *storeHandler) Block(c *fiber.Ctx) error {
	span, ctx := opentracing.StartSpanFromContext(c.Context(), "StoreHandler.Block")
	defer span.Finish()
	blockRequests.Inc()

	id := c.Params("id")

	err := h.validate.VarCtx(ctx, id, "uuid4")
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.VarCtx: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	err = h.storeUsecase.Block(ctx, id)
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Block: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	successRequests.Inc()
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Disable stores
// @Description Disable a stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {array} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id}/disable [patch]
func (h *storeHandler) Disable(c *fiber.Ctx) error {
	span, ctx := opentracing.StartSpanFromContext(c.Context(), "StoreHandler.Disable")
	defer span.Finish()
	disableRequests.Inc()

	id := c.Params("id")

	err := h.validate.VarCtx(ctx, id, "uuid4")
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.VarCtx: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	err = h.storeUsecase.Disable(ctx, id)
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Disable: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	successRequests.Inc()
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
// @Failure 400 {array} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id} [delete]
func (h *storeHandler) Delete(c *fiber.Ctx) error {
	span, ctx := opentracing.StartSpanFromContext(c.Context(), "StoreHandler.Delete")
	defer span.Finish()
	deleteRequests.Inc()

	id := c.Params("id")

	err := h.validate.VarCtx(ctx, id, "uuid4")
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.VarCtx %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	err = h.storeUsecase.Delete(ctx, id)
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Delete: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	successRequests.Inc()
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Update store
// @Description Uptate a stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Param category body domain.UpdateStoreRequest true "Create store"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {array} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id} [patch]
func (h *storeHandler) Update(c *fiber.Ctx) error {
	span, ctx := opentracing.StartSpanFromContext(c.Context(), "StoreHandler.Update")
	defer span.Finish()
	updateRequests.Inc()

	ur := new(domain.UpdateStoreRequest)
	id := c.Params("id")

	if err := c.BodyParser(ur); err != nil {
		log.
			WithContext(ctx).
			Errorf("c.BodyParser: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	ur.ID = id
	if err := h.validate.StructCtx(ctx, ur); err != nil {
		log.
			WithContext(ctx).
			Errorf("validate.StructCtx: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}
	err := h.storeUsecase.Update(ctx, ur)
	if err != nil {
		log.
			WithContext(ctx).
			Errorf("storeUsecase.Update: %v", err)
		errorRequests.Inc()
		return errorHandler(c, err)
	}

	successRequests.Inc()
	return c.SendStatus(fiber.StatusNoContent)
}
