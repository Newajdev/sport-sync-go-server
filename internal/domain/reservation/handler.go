package reservation

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/domain/zone"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func getCurrentUserID(c *echo.Context) (uint, bool) {
	userID, ok := c.Get("user_id").(uint)
	return userID, ok
}

func reservationErrorResponse(c *echo.Context, err error) error {
	if errors.Is(err, ErrReservationNotFound) {
		return c.JSON(http.StatusNotFound, httpresponse.Fail(
			"Reservation not found",
			err.Error(),
		))
	}

	if errors.Is(err, zone.ErrZoneNotFound) {
		return c.JSON(http.StatusNotFound, httpresponse.Fail(
			"Parking zone not found",
			err.Error(),
		))
	}

	if errors.Is(err, ErrZoneFull) {
		return c.JSON(http.StatusConflict, httpresponse.Fail(
			"Zone is at full capacity",
			err.Error(),
		))
	}

	if errors.Is(err, ErrReservationAlreadyCancelled) {
		return c.JSON(http.StatusConflict, httpresponse.Fail(
			"Reservation is already cancelled",
			err.Error(),
		))
	}

	if errors.Is(err, ErrForbiddenReservationAccess) {
		return c.JSON(http.StatusForbidden, httpresponse.Fail(
			"You do not own this reservation",
			err.Error(),
		))
	}

	return c.JSON(http.StatusInternalServerError, httpresponse.Fail(
		"Something went wrong",
		err.Error(),
	))
}

func (h *handler) Create(c *echo.Context) error {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Fail(
			"Unauthorized",
			"Missing or invalid user context",
		))
	}

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

	response, err := h.service.CreateReservation(userID, req)
	if err != nil {
		return reservationErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, httpresponse.OK("Reservation confirmed successfully", response))
}

func (h *handler) GetMyReservations(c *echo.Context) error {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Fail(
			"Unauthorized",
			"Missing or invalid user context",
		))
	}

	reservations, err := h.service.GetMyReservations(userID)
	if err != nil {
		return reservationErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, httpresponse.OK("My reservations retrieved successfully", reservations))
}

func (h *handler) Cancel(c *echo.Context) error {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Fail(
			"Unauthorized",
			"Missing or invalid user context",
		))
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Fail(
			"Invalid reservation id",
			err.Error(),
		))
	}

	if err := h.service.CancelReservation(userID, uint(id)); err != nil {
		return reservationErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, httpresponse.OK("Reservation cancelled successfully", nil))
}

func (h *handler) GetAll(c *echo.Context) error {
	reservations, err := h.service.GetAllReservations()
	if err != nil {
		return reservationErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, httpresponse.OK("Reservations retrieved successfully", reservations))
}
