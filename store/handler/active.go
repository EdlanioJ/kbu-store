package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type ActiveHandler struct {
	activeUsecase domain.ActivateStoreUsecase
}

func NewActiveHandler(u domain.ActivateStoreUsecase) *ActiveHandler {
	return &ActiveHandler{
		activeUsecase: u,
	}
}

// @Summary Active Store
// @Description Set Store status to active
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 404 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/{id}/active [options]
func (s *ActiveHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}

	err := s.activeUsecase.Exec(ctx, id)

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
