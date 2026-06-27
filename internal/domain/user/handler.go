package user

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) Register(c echo.Context) error {
	var req dto.RegisterRequest

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

	response, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, ErrorAlreadyExist) {
			return c.JSON(http.StatusConflict, httpresponse.Fail(
				"Failed to register user",
				err.Error(),
			))
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Fail(
			"Failed to register user",
			err.Error(),
		))
	}

	return c.JSON(http.StatusCreated, httpresponse.OK("User registered successfully", response))
}

func (h *handler) Login(c echo.Context) error {
	var req dto.LoginRequest

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

	response, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, httpresponse.Fail(
				"Login failed",
				err.Error(),
			))
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Fail(
			"Failed to login",
			err.Error(),
		))
	}

	return c.JSON(http.StatusOK, httpresponse.OK("Login successful", response))
}
