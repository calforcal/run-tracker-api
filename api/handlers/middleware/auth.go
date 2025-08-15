package middleware

import (
	"fmt"
	"net/http"
	"run-tracker-api/internal/auth"
	"run-tracker-api/internal/config"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	AuthMiddleware struct {
		config  *config.Config
		service *auth.AuthService
	}
)

func NewAuthMiddleware(cfg *config.Config, s *auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		config:  cfg,
		service: s,
	}
}

func (m *AuthMiddleware) RunAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing or invalid token"})
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := m.service.ParseJWT(tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
			}

			// Set UUID in context for downstream use
			c.Set("uuid", claims.UUID)
			return next(c)
		}
	}
}
