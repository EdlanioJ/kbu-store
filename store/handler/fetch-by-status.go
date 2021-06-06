package handler

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/gofiber/fiber/v2"
)

type FetchByStatusHandler struct {
	fetchByStatusUsecase domain.FetchStoreByStatusUsecase
}

func NewFetchByStatusHandler(u domain.FetchStoreByStatusUsecase) *FetchByStatusHandler {
	return &FetchByStatusHandler{
		fetchByStatusUsecase: u,
	}
}

// @Summary Get all stores by status
// @Description Get all stores by status
// @Tags stores
// @Accept json
// @Produce json
// @Param status path string true "Status" Enums(active,pending,block,disable)
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Success 200 {array} domain.Store
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/status/{status} [get]
func (s *FetchByStatusHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	sort, _ := url.QueryUnescape(c.Query("sort"))

	status := c.Params("status")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	list, total, err := s.fetchByStatusUsecase.Exec(ctx, status, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}
