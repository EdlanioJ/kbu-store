package handler

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/EdlanioJ/kbu-store/domain"
	"github.com/EdlanioJ/kbu-store/dto"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type FetchByLocationAndStatusHandler struct {
	fetchByLocationUsecase domain.FetchStoreByCloseLocationUsecase
}

func NewFetchByLocationAndStatusHandler(usecase domain.FetchStoreByCloseLocationUsecase) *FetchByLocationAndStatusHandler {
	return &FetchByLocationAndStatusHandler{
		fetchByLocationUsecase: usecase,
	}
}

// @Summary Get all stores by location and status
// @Description Get all stores by location and status
// @Tags stores
// @Accept json
// @Produce json
// @Param status path string true "Status value" Enums(active,pending,block,disable)
// @Param page query int false "List Page"
// @Param limit query int false "List Limit"
// @Param sort query string false "List Sort"
// @Param distance query int true "Distance in KM"
// @Param location path string true "Location @lat,lng"
// @Success 200 {array} domain.Store
// @Failure 500 {object} dto.ErrorResponse{}
// @Router /stores/location/{location}/status/{status} [get]
func (s *FetchByLocationAndStatusHandler) Handler(c *fiber.Ctx) error {
	ctx := c.Context()
	sort, _ := url.QueryUnescape(c.Query("sort"))

	status := c.Params("status")
	distance, _ := strconv.Atoi(c.Query("distance"))
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	locationUnescape, _ := url.PathUnescape(c.Params("location"))
	locationPath := strings.Replace(locationUnescape, "@", "", -1)

	location := strings.Split(locationPath, ",")

	if !govalidator.IsLatitude(strings.TrimSpace(location[0])) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "must be a valid latitude",
		})
	}

	if !govalidator.IsLongitude(strings.TrimSpace(location[1])) {
		logrus.Info(location)
		logrus.Info(locationUnescape)
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "must be a valid longitude",
		})
	}

	lat, _ := strconv.ParseFloat(strings.TrimSpace(location[0]), 8)
	lng, _ := strconv.ParseFloat(strings.TrimSpace(location[1]), 8)

	stores, total, err := s.fetchByLocationUsecase.Exec(ctx, lat, lng, distance, status, limit, page, sort)

	c.Response().Header.Add("X-total", fmt.Sprint(total))

	if err != nil {
		return c.Status(getStatusCode(err)).JSON(dto.ErrorResponse{
			Message: err.Error(),
		})
	}
	return c.JSON(stores)
}
