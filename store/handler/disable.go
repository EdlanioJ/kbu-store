package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type DisableHandler struct {
	disableUsecase domain.DisableStoreUsecase
}

func NewDisableHandler(u domain.DisableStoreUsecase) *DisableHandler {
	return &DisableHandler{
		disableUsecase: u,
	}
}

// @Summary Disable Store
// @Description Set Store status to disable
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 404 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/{id}/disable [options]
func (s *DisableHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}

	err := s.disableUsecase.Exec(ctx, id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
