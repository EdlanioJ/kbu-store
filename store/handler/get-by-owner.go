package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type GetByOwnerHandler struct {
	getByOwnerUsecase domain.GetStoreByOwnerUsecase
}

func NewGetByOwnerHandler(u domain.GetStoreByOwnerUsecase) *GetByOwnerHandler {
	return &GetByOwnerHandler{
		getByOwnerUsecase: u,
	}
}

// @Summary Get store by owner
// @Description Get store by owner
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Param owner path string true "owner ID"
// @Success 200 {object} domain.Store
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 404 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/{id}/owner/{owner} [get]
func (s *GetByOwnerHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	owner := c.Params("owner")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}

	if !govalidator.IsUUIDv4(owner) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "owner must be a valid uuidv4",
			})
	}

	res, err := s.getByOwnerUsecase.Exec(ctx, id, owner)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(res)
}
