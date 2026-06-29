package middlewares

import (
	"net/http"
	"spotsync/internal/domain/user"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role, ok := c.Get("user_role").(string)
			if !ok || user.Role(role) != user.RoleAdmin {
				return c.JSON(http.StatusForbidden, httpresponse.Fail(
					"Forbidden",
					"Admin access required",
				))
			}

			return next(c)
		}
	}
}
