package middleware

import (
	"net/http"
	"strings"

	"run-tracker-api/internal/ports"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	service ports.AuthService
}

func NewAuthMiddleware(service ports.AuthService) *AuthMiddleware {
	return &AuthMiddleware{service: service}
}

func (m *AuthMiddleware) RunAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing or invalid token"})
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := m.service.ParseToken(c.Request().Context(), tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
			}

			// Set UUID in context for downstream use
			c.Set("uuid", claims.UUID)
			return next(c)
		}
	}
}
