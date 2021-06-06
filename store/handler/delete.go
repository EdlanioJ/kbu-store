package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type DeleteHandler struct {
	deleteUsecase domain.DeleteStoreUsecase
}

func NewDeleteHandler(u domain.DeleteStoreUsecase) *DeleteHandler {
	return &DeleteHandler{
		deleteUsecase: u,
	}
}

// @Summary Delete Store
// @Description Delete store
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 404 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/{id} [delete]
func (s *DeleteHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}

	err := s.deleteUsecase.Exec(ctx, id)

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
