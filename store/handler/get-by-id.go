package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type GetByIDHandler struct {
	getByIDUsecase domain.GetStoreByIDUsecase
}

func NewGetByIDHandler(usecase domain.GetStoreByIDUsecase) *GetByIDHandler {
	return &GetByIDHandler{
		getByIDUsecase: usecase,
	}
}

// @Summary Get store by id
// @Description Get store by id
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 200 {object} domain.Store
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 404 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/{id} [get]
func (s *GetByIDHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	if !govalidator.IsUUIDv4(id) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Message: "id must be a valid uuidv4",
			})
	}

	res, err := s.getByIDUsecase.Exec(ctx, id)

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(res)
}
