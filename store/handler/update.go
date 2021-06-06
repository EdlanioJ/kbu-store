package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type UpdateHandler struct {
	uptateUsecase domain.UpdateStoreUsecase
}

func NewUpdateHandler(u domain.UpdateStoreUsecase) *UpdateHandler {
	return &UpdateHandler{
		uptateUsecase: u,
	}
}

// @Summary Get store by id
// @Description Get store by id
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Param store body dto.UpdateStoreRequest true "Update store request body"
// @Success 200 {object} domain.Store
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 422 {object} dto.ErrorResponse{}
// @Failure 404 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/{id} [patch]
func (s *UpdateHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	reqBody := new(dto.UpdateStoreRequest)
	id := c.Params("id")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}

	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	store := reqBody.Parser()
	store.ID = id

	err := s.uptateUsecase.Exec(ctx, store)

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)

}
