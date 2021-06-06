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

type FetchByOwnerHandler struct {
	fetchByOwnerUsecase domain.FetchStoreByOwnerUsecase
}

func NewFetchByOwnerHandler(u domain.FetchStoreByOwnerUsecase) *FetchByOwnerHandler {
	return &FetchByOwnerHandler{
		fetchByOwnerUsecase: u,
	}
}

// @Summary Get all stores by owner
// @Description Get all stores by owner
// @Tags stores
// @Accept json
// @Produce json
// @Param owner path string true "owner ID"
// @Param page query int false "List Page"
// @Param limit query int false "List Limit"
// @Param sort query string false "List Sort"
// @Success 200 {array} domain.Store
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/owner/{owner} [get]
func (s *FetchByOwnerHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	sort, _ := url.QueryUnescape(c.Query("sort"))

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	externalID := c.Params("owner")

	if !govalidator.IsUUIDv4(externalID) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "owner must be a valid uuidv4",
			})
	}

	list, total, err := s.fetchByOwnerUsecase.Exec(ctx, externalID, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}
