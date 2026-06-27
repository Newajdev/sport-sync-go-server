package zone

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

	api := e.Group("/api/v1/zones")

	api.POST("", h.Create, middlewares.AuthMiddleware(jwtService), middlewares.AdminMiddleware())
	api.GET("", h.GetAll)
	api.GET("/:id", h.GetByID)
}
