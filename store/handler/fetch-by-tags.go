package handler

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/gofiber/fiber/v2"
)

type FetchByTagsHandler struct {
	fetchByTagsUsecase domain.FetchStoreByTagsUsecase
}

func NewFetchByTagsHandler(u domain.FetchStoreByTagsUsecase) *FetchByTagsHandler {
	return &FetchByTagsHandler{
		fetchByTagsUsecase: u,
	}
}

// @Summary Get all stores by tags
// @Description Get all stores by tags
// @Tags stores
// @Accept json
// @Produce json
// @Param tags query []string true "tags"
// @Param page query int false "List Page"
// @Param limit query int false "List Limit"
// @Param sort query string false "List Sort"
// @Success 200 {array} domain.Store
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/tags/ [get]
func (s *FetchByTagsHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()

	sort, _ := url.QueryUnescape(c.Query("sort"))

	tags := strings.Split(c.Query("tags"), ",")

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	list, total, err := s.fetchByTagsUsecase.Exec(ctx, tags, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}
