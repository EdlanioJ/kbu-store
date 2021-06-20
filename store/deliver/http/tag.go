package http

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type tagHandler struct {
	tagUsecase domain.TagUsecase
}

func NewTagRoutes(r fiber.Router, us domain.TagUsecase) {
	handler := &tagHandler{
		tagUsecase: us,
	}
	tr := r.Group("/tags")

	tr.Get("/", handler.getAll)
	tr.Get("/category/:category", handler.getAllByCategory)
}

// @Summary Get all tags
// @Description Get all tags
// @Tags tags
// @Accept json
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param sort query string false "Sort" default(count DESC)
// @Success 200 {array} domain.Tag
// @Failure 500 {object} ErrorResponse
// @Router /tags [get]
func (h *tagHandler) getAll(c *fiber.Ctx) error {
	ctx := c.Context()

	sort, _ := url.QueryUnescape(c.Query("sort"))
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	list, total, err := h.tagUsecase.GetAll(ctx, sort, page, limit)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}

// @Summary Get all tags by category
// @Description Get all tags by category
// @Tags tags
// @Accept json
// @Produce json
// @Param category path string true "category ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param sort query string false "Sort" default(count DESC)
// @Success 200 {array} domain.Tag
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tags/category/{category} [get]
func (h *tagHandler) getAllByCategory(c *fiber.Ctx) error {
	ctx := c.Context()
	sort, _ := url.QueryUnescape(c.Query("sort"))

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	categoryID := c.Params("category")
	if !govalidator.IsUUIDv4(categoryID) {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Message: "category must be a valid uuidv4",
			})
	}
	list, total, err := h.tagUsecase.GetAllByCategory(ctx, categoryID, sort, limit, page)
	c.Response().Header.Add("X-total", fmt.Sprint(total))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(list)
}
