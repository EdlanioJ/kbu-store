package handler

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/validators"
	"github.com/gofiber/fiber/v2"
)

type storeHandler struct {
	storeUsecase domain.StoreUsecase
}

func NewStoreRoute(r fiber.Router, us domain.StoreUsecase) {
	handler := &storeHandler{
		storeUsecase: us,
	}
	sr := r.Group("/stores")

	sr.Post("/", handler.create)
	sr.Get("/", handler.getAll)
	sr.Get("/category/:category", handler.getAllByCategory)
	sr.Get("/location/:location/status/:status", handler.getAllByCloseLocation)
	sr.Get("/owner/:owner", handler.getAllByOwner)
	sr.Get("/status/:status", handler.getAllByStatus)
	sr.Get("/tags/", handler.getAllByTags)
	sr.Get("/:id", handler.getById)
	sr.Get("/:id/owner/:owner", handler.getByIdAndOwner)

	sr.Patch("/:id", handler.update)
	sr.Patch("/:id/activate", handler.active)
	sr.Patch("/:id/block", handler.block)
	sr.Patch("/:id/disable", handler.disable)

	sr.Delete("/:id", handler.delete)

}

// @Summary Create  a new store
// @Description Create store
// @Tags stores
// @Accept json
// @Produce json
// @Param category body CreateStoreRequest true "Create store"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /stores [post]
func (h *storeHandler) create(c *fiber.Ctx) error {
	ctx := c.Context()
	cr := new(CreateStoreRequest)
	if err := c.BodyParser(cr); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	err := h.storeUsecase.Create(ctx, cr.Name, cr.Description, cr.CategoryID, cr.ExternalID, cr.Tags, cr.Lat, cr.Lng)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}

// @Summary Get all stores
// @Description Get a list of stores
// @Tags stores
// @Accept json
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param sort query string false "Sort" default(created_at DESC)
// @Success 200 {array} domain.Store
// @Failure 500 {object} ErrorResponse
// @Router /stores [get]
func (h *storeHandler) getAll(c *fiber.Ctx) error {
	sort := c.Query("sort")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	ctx := c.Context()

	list, total, err := h.storeUsecase.GetAll(ctx, sort, limit, page)
	c.Response().Header.Add("X-total", fmt.Sprint(total))
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(list)
}

// @Summary Get all stores by category
// @Description Get a list of stores by category
// @Tags stores
// @Accept json
// @Produce json
// @Param category path string true "category ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param sort query string false "Sort" default(created_at DESC)
// @Success 200 {array} domain.Store
// @Failure 500 {object} ErrorResponse
// @Router /stores/category/{category} [get]
func (h *storeHandler) getAllByCategory(c *fiber.Ctx) error {
	ctx := c.Context()
	sort, _ := url.PathUnescape(c.Query("sort"))
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	categoryID := c.Params("category")
	err := validators.ValidateUUIDV4("category", categoryID)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(
			ErrorResponse{
				Message: err.Error(),
			})
	}

	list, total, err := h.storeUsecase.GetAllByCategory(ctx, categoryID, sort, limit, page)
	c.Response().Header.Add("X-total", fmt.Sprint(total))
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(list)
}

// @Summary Get all stores location and status
// @Description Get a list of stores location and status
// @Tags stores
// @Accept json
// @Produce json
// @Param location path string true "@lat,lng"
// @Param status path string true "Status" Enums(active,pending,block,disable)
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param sort query string false "Sort" default(created_at DESC)
// @Success 200 {array} domain.Store
// @Failure 500 {object} ErrorResponse
// @Router /stores/location/{location}/status/{status} [get]
func (h *storeHandler) getAllByCloseLocation(c *fiber.Ctx) error {
	ctx := c.Context()
	sort, _ := url.QueryUnescape(c.Query("sort"))

	status := c.Params("status")
	distance, _ := strconv.Atoi(c.Query("distance"))
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	locationUnescape, _ := url.PathUnescape(c.Params("location"))
	locationPath := strings.Replace(locationUnescape, "@", "", -1)

	location := strings.Split(locationPath, ",")

	lat, _ := strconv.ParseFloat(strings.TrimSpace(location[0]), 64)
	lng, _ := strconv.ParseFloat(strings.TrimSpace(location[1]), 64)

	err := validators.ValidateLatitude(lat)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}

	err = validators.ValidateLongitude(lng)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	stores, total, err := h.storeUsecase.GetAllByByCloseLocation(ctx, lat, lng, distance, status, limit, page, sort)
	c.Response().Header.Add("X-total", fmt.Sprint(total))
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(stores)
}

// @Summary Get all stores by owner
// @Description Get a list of stores owner
// @Tags stores
// @Accept json
// @Produce json
// @Param owner path string true "user ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param sort query string false "Sort" default(created_at DESC)
// @Success 200 {array} domain.Store
// @Failure 500 {object} ErrorResponse
// @Router /stores/owner/{owner} [get]
func (h *storeHandler) getAllByOwner(c *fiber.Ctx) error {
	ctx := c.Context()
	sort, _ := url.QueryUnescape(c.Query("sort"))
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	externalID := c.Params("owner")

	err := validators.ValidateUUIDV4("owner", externalID)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(
			ErrorResponse{
				Message: err.Error(),
			})
	}

	list, total, err := h.storeUsecase.GetAllByOwner(ctx, externalID, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(list)
}

// @Summary Get all stores by status
// @Description Get a list of stores by status
// @Tags stores
// @Accept json
// @Produce json
// @Param status path string true "Status" Enums(active,pending,block,disable)
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param sort query string false "Sort" default(created_at DESC)
// @Success 200 {array} domain.Store
// @Failure 500 {object} ErrorResponse
// @Router /stores/status/{status} [get]
func (h *storeHandler) getAllByStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	sort, _ := url.QueryUnescape(c.Query("sort"))
	status := c.Params("status")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	list, total, err := h.storeUsecase.GetAllByStatus(ctx, status, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(list)
}

// @Summary Get all stores
// @Description Get a list of stores
// @Tags stores
// @Accept json
// @Produce json
// @Param tags query []string true "Tags"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param sort query string false "Sort" default(created_at DESC)
// @Success 200 {array} domain.Store
// @Failure 500 {object} ErrorResponse
// @Router /stores/tags/ [get]
func (h *storeHandler) getAllByTags(c *fiber.Ctx) error {
	ctx := c.Context()

	sort, _ := url.QueryUnescape(c.Query("sort"))

	tags := strings.Split(c.Query("tags"), ",")

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	list, total, err := h.storeUsecase.GetAllByTags(ctx, tags, sort, limit, page)

	c.Response().Header.Add("X-total", fmt.Sprint(total))
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(list)
}

// @Summary Get stores by id
// @Description Get one stores by id
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 200 {object} domain.Store
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id} [get]
func (h *storeHandler) getById(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(
			ErrorResponse{
				Message: err.Error(),
			})
	}

	res, err := h.storeUsecase.GetById(ctx, id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(res)
}

// @Summary Get stores by id and owner
// @Description Get one stores by id and owner
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Param owner path string true "user ID"
// @Success 200 {object} domain.Store
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id}/owner/{owner} [get]
func (h *storeHandler) getByIdAndOwner(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	owner := c.Params("owner")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(
			ErrorResponse{
				Message: err.Error(),
			})
	}

	err = validators.ValidateUUIDV4("owner", owner)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(
			ErrorResponse{
				Message: err.Error(),
			})
	}

	res, err := h.storeUsecase.GetByIdAndOwner(ctx, id, owner)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(res)
}

// @Summary Activate stores
// @Description Activate one stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id}/activate [patch]
func (h *storeHandler) active(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(
			ErrorResponse{
				Message: err.Error(),
			})
	}

	err = h.storeUsecase.Active(ctx, id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Block stores
// @Description Block one stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id}/block [patch]
func (h *storeHandler) block(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(
			ErrorResponse{
				Message: err.Error(),
			})
	}

	err = h.storeUsecase.Block(ctx, id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Disable stores
// @Description Disable one stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id}/disable [patch]
func (h *storeHandler) disable(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(
			ErrorResponse{
				Message: err.Error(),
			})
	}

	err = h.storeUsecase.Disable(ctx, id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Delete stores
// @Description Delete one stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/{id} [delete]
func (h *storeHandler) delete(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(
			ErrorResponse{
				Message: err.Error(),
			})
	}

	err = h.storeUsecase.Delete(ctx, id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Activate stores
// @Description Activate one stores
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "store ID"
// @Param category body UpdateStoreRequest true "Create store"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /stores/active/{id} [patch]
func (h *storeHandler) update(c *fiber.Ctx) error {
	ctx := c.Context()
	reqBody := new(UpdateStoreRequest)
	id := c.Params("id")

	err := validators.ValidateUUIDV4("id", id)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(
			ErrorResponse{
				Message: err.Error(),
			})
	}

	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}

	store := reqBody.ToDomainStore()
	store.ID = id

	err = h.storeUsecase.Update(ctx, store)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
