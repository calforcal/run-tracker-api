package home

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
}

func New() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Home(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
