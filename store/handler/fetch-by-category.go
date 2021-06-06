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

type FetchByTypeHandler struct {
	fetchByTypeUsecase domain.FetchStoreByTypeUsecase
}

func NewFetchByTypeHandler(u domain.FetchStoreByTypeUsecase) *FetchByTypeHandler {
	return &FetchByTypeHandler{
		fetchByTypeUsecase: u,
	}
}

// @Summary Get all stores by category
// @Description Get all stores by category
// @Tags stores
// @Accept json
// @Produce json
// @Param category path string true "category ID"
// @Param page query int false "List Page"
// @Param limit query int false "List Limit"
// @Param sort query string false "List Sort"
// @Success 200 {array} domain.Store
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/category/{category} [get]
func (s *FetchByTypeHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	sort, _ := url.PathUnescape(c.Query("sort"))

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	categoryID := c.Params("category")

	if !govalidator.IsUUIDv4(categoryID) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "category must be a valid uuidv4",
			})
	}

	list, total, err := s.fetchByTypeUsecase.Exec(ctx, categoryID, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}
