package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type GetByStatusHandler struct {
	getByStatusUsecase domain.GetCategoryByStautsUsecase
}

func NewGetByStatusHandler(usecase domain.GetCategoryByStautsUsecase) *GetByStatusHandler {
	return &GetByStatusHandler{
		getByStatusUsecase: usecase,
	}
}

// @Summary Get all categories by status
// @Description Get all categories by status
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "category ID"
// @Param status path string true "Status" Enums(active,pending,block,disable)
// @Success 200 {object} domain.Category
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 404 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /categories/{id}/{status} [get]
func (s *GetByStatusHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	status := c.Params("status")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}

	res, err := s.getByStatusUsecase.Exec(ctx, id, status)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(res)
}
