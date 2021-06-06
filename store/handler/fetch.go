package handler

import (
	"fmt"
	"strconv"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/gofiber/fiber/v2"
)

type FetchHandler struct {
	fetchUsecase domain.FetchStoreUsecase
}

func NewFetchHandler(usecase domain.FetchStoreUsecase) *FetchHandler {
	return &FetchHandler{
		fetchUsecase: usecase,
	}
}

// @Summary Get all stores
// @Description Get all stores
// @Tags stores
// @Accept json
// @Produce json
// @Param page query int false "List Page"
// @Param limit query int false "List Limit"
// @Param sort query string false "List Sort"
// @Success 200 {array} domain.Category
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores [get]
func (s *FetchHandler) Handler(c *fiber.Ctx) error {
	sort := c.Query("sort")

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	ctx := c.Context()

	list, total, err := s.fetchUsecase.Exec(ctx, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}
