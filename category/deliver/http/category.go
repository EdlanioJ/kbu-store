package http

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type categoryHandler struct {
	categoryUsecase domain.CategoryUsecase
}

func NewCategoryRoutes(r fiber.Router, us domain.CategoryUsecase) {
	handler := &categoryHandler{
		categoryUsecase: us,
	}
	cr := r.Group("/categories")

	cr.Post("/", handler.create)
	cr.Get("/", handler.getAll)
	cr.Get("/:id", handler.getById)
	cr.Get("/:id/status/:status", handler.getByIdAndStatus)
	cr.Get("/status/:status", handler.getAllByStatus)

	cr.Patch("/:id/activate", handler.activate)
	cr.Patch("/:id/disable", handler.disable)
}

// @Summary Create a new Category
// @Description Create category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CreateCategoryRequest true "Create category"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /categories [post]
func (h *categoryHandler) create(c *fiber.Ctx) error {
	ctx := c.Context()
	cr := new(CreateCategoryRequest)

	if err := c.BodyParser(cr); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	err := h.categoryUsecase.Create(ctx, cr.Name)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}

// @Summary Get categories by id
// @Description Get all categories by status
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "category ID"
// @Success 200 {object} domain.Category
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [get]
func (h *categoryHandler) getById(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}
	res, err := h.categoryUsecase.GetById(ctx, id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(res)
}

// @Summary Get all categories by status
// @Description Get all categories by status
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "category ID"
// @Param status path string true "Status" Enums(active,pending,block,disable)
// @Success 200 {object} domain.Category
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id}/status/{status} [get]
func (h *categoryHandler) getByIdAndStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	status := c.Params("status")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}

	res, err := h.categoryUsecase.GetByIdAndStatus(ctx, id, status)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(res)
}

// @Summary Get all categories
// @Description Get all categories
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "List Page"
// @Param limit query int false "List Limit"
// @Param sort query string false "List Sort"
// @Success 200 {array} domain.Category
// @Failure 500 {object} ErrorResponse
// @Router /categories [get]
func (h *categoryHandler) getAll(c *fiber.Ctx) error {
	sort, _ := url.QueryUnescape(c.Query("sort"))
	p := c.Query("page")
	l := c.Query("limit")

	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)
	ctx := c.Context()

	list, total, err := h.categoryUsecase.GetAll(ctx, sort, limit, page)
	c.Response().Header.Add("X-total", fmt.Sprint(total))
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(list)
}

// @Summary Get all categories by status
// @Description Get all categories by status
// @Tags categories
// @Accept json
// @Produce json
// @Param status path string true "Status" Enums(active,pending,block,disable)
// @Param page query int false "List Page"
// @Param limit query int false "List Limit"
// @Param sort query string false "List Sort"
// @Success 200 {array} domain.Category
// @Failure 500 {object} ErrorResponse
// @Router /categories/status/{status} [get]
func (h *categoryHandler) getAllByStatus(c *fiber.Ctx) error {
	sort, _ := url.QueryUnescape(c.Query("sort"))

	p := c.Query("page")
	l := c.Query("limit")

	status := c.Params("status")
	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)
	ctx := c.Context()

	list, total, err := h.categoryUsecase.GetAllByStatus(ctx, status, sort, limit, page)
	c.Response().Header.Add("X-total", fmt.Sprint(total))
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(list)
}

// @Summary Activate categories
// @Description Activate one category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "category ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /categories/activate/{id} [patch]
func (h *categoryHandler) activate(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}
	err := h.categoryUsecase.Activate(ctx, id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Disable categories
// @Description Disable one category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "category ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /categories/disable/{id} [patch]
func (h *categoryHandler) disable(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}
	err := h.categoryUsecase.Disable(ctx, id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
