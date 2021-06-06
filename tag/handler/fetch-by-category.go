package handler

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type FetchTagsByCategory struct {
	fetchTagsByCategoryUsecase domain.FetchTagsByCategoryUsecase
}

func NewFetchTagsByCategory(usecase domain.FetchTagsByCategoryUsecase) *FetchTagsByCategory {
	return &FetchTagsByCategory{
		fetchTagsByCategoryUsecase: usecase,
	}
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
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /tags/category/{category} [get]
func (h *FetchTagsByCategory) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	sort, _ := url.QueryUnescape(c.Query("sort"))

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	categoryID := c.Params("category")

	if !govalidator.IsUUIDv4(categoryID) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "category must be a valid uuidv4",
			})
	}

	list, total, err := h.fetchTagsByCategoryUsecase.Exec(ctx, categoryID, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}
