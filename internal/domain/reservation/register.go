package reservation

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	jwtService := auth.NewJWTService(cfg.JwtSecret)

	api := e.Group("/api/v1/reservations", middlewares.AuthMiddleware(jwtService))

	api.POST("", h.Create)
	api.GET("/my-reservations", h.GetMyReservations)
	api.DELETE("/:id", h.Cancel)
	api.GET("", h.GetAll, middlewares.AdminMiddleware())
}
