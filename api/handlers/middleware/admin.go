package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AdminMiddleware gates developer-only endpoints (e.g. webhook subscription
// management) behind a single static bearer token, distinct from per-user
// JWTs and from any third-party webhook handshake secret.
type AdminMiddleware struct {
	token string
}

func NewAdminMiddleware(token string) *AdminMiddleware {
	return &AdminMiddleware{token: token}
}

func (m *AdminMiddleware) RunAdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			tokenStr, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok || m.token == "" || subtle.ConstantTimeCompare([]byte(tokenStr), []byte(m.token)) != 1 {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing or invalid admin token"})
			}
			return next(c)
		}
	}
}
