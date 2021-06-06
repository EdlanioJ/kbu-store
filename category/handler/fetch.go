package handler

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/gofiber/fiber/v2"
)

type FetchHandler struct {
	fetchUsecase domain.FetchCategoryUsecase
}

func NewFetchHandler(usecase domain.FetchCategoryUsecase) *FetchHandler {
	return &FetchHandler{
		fetchUsecase: usecase,
	}
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
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /categories [get]
func (st *FetchHandler) Handler(c *fiber.Ctx) error {
	sort, _ := url.QueryUnescape(c.Query("sort"))

	p := c.Query("page")
	l := c.Query("limit")

	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)
	ctx := c.Context()

	list, total, err := st.fetchUsecase.Exec(ctx, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}
