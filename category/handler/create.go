package handler

import (
	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/gofiber/fiber/v2"
)

type CreateHandler struct {
	createUsecase domain.CreateCategoryUsecase
}

func NewCreateHandler(u domain.CreateCategoryUsecase) *CreateHandler {
	return &CreateHandler{
		createUsecase: u,
	}
}

// @Summary Create  a new Category
// @Description Create category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body dto.CreateCategoryRequest true "Create category"
// @Success 201 {string} string "Created"
// @Failure 400 {object} dto.ErrorResponse{}
// @Failure 500 {object} dto.ErrorResponse{}
// @Failure 422 {object} dto.ErrorResponse{}
// @Router /categories [post]
func (st *CreateHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	cr := new(dto.CreateCategoryRequest)

	if err := c.BodyParser(cr); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	err := st.createUsecase.Add(ctx, cr.Name)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}
