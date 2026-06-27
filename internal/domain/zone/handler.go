package zone

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync/internal/domain/zone/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) Create(c *echo.Context) error {
	var req dto.CreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Fail(
			"Invalid request payload",
			err.Error(),
		))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Fail(
			"Validation failed",
			err.Error(),
		))
	}

	response, err := h.service.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Fail(
			"Failed to create parking zone",
			err.Error(),
		))
	}

	return c.JSON(http.StatusCreated, httpresponse.OK("Parking zone created successfully", response))
}

func (h *handler) GetAll(c *echo.Context) error {
	zones, err := h.service.GetZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Fail(
			"Failed to retrieve parking zones",
			err.Error(),
		))
	}

	return c.JSON(http.StatusOK, httpresponse.OK("Parking zones retrieved successfully", zones))
}

func (h *handler) GetByID(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Fail(
			"Invalid zone id",
			err.Error(),
		))
	}

	response, err := h.service.GetZoneByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Fail(
				"Parking zone not found",
				err.Error(),
			))
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Fail(
			"Failed to retrieve parking zone",
			err.Error(),
		))
	}

	return c.JSON(http.StatusOK, httpresponse.OK("Parking zone retrieved successfully", response))
}
