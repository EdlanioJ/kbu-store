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
	fetchByStatusUsecase domain.FetchCategoryByStatusUsecase
}

func NewFetchByStatusHandler(usecase domain.FetchCategoryByStatusUsecase) *FetchByStatusHandler {
	return &FetchByStatusHandler{
		fetchByStatusUsecase: usecase,
	}
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
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /categories/status/{status} [get]
func (st *FetchByStatusHandler) Handler(c *fiber.Ctx) error {
	sort, _ := url.QueryUnescape(c.Query("sort"))

	p := c.Query("page")
	l := c.Query("limit")

	status := c.Params("status")
	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)
	ctx := c.Context()

	list, total, err := st.fetchByStatusUsecase.Exec(ctx, status, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}
