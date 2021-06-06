package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/gofiber/fiber/v2"
)

type CreateHandler struct {
	createUsecase domain.CreateStoreUsecase
}

func NewCreateHandler(u domain.CreateStoreUsecase) *CreateHandler {
	return &CreateHandler{
		createUsecase: u,
	}
}

// @Summary Create  a new Store
// @Description Create store
// @Tags stores
// @Accept json
// @Produce json
// @Param store body dto.CreateStoreRequest true "Create store request"
// @Success 201 {string} string "Created"
// @Failure 422 {object} dto.ErrorResponse{}
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores [post]
func (s *CreateHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	cr := new(dto.CreateStoreRequest)

	if err := c.BodyParser(cr); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	err := s.createUsecase.Add(ctx, cr.Name, cr.Description, cr.CategoryID, cr.ExternalID, cr.Tags, cr.Lat, cr.Lng)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}
