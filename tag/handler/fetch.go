package handler

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/gofiber/fiber/v2"
)

type FetchTags struct {
	fetchTagsUsecase domain.FetchTagsUsecase
}

func NewFetchTags(usecase domain.FetchTagsUsecase) *FetchTags {
	return &FetchTags{
		fetchTagsUsecase: usecase,
	}
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
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /tags [get]
func (h *FetchTags) Handler(c *fiber.Ctx) error {
	ctx := c.Context()

	sort, _ := url.QueryUnescape(c.Query("sort"))
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	list, total, err := h.fetchTagsUsecase.Exec(ctx, sort, page, limit)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}
