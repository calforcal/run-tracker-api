package webhooks

import (
	"net/http"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/users"

	"github.com/labstack/echo/v4"
)

type (
	WebhookHandler struct {
		cfg         *config.Config
		userService *users.UserService
	}
)

func New(cfg *config.Config, userService *users.UserService) WebhookHandler {
	return WebhookHandler{
		cfg:         cfg,
		userService: userService,
	}
}

func (h *WebhookHandler) ProcessAthleteUpload(c echo.Context) error {
	
	
	return c.JSON(http.StatusOK, map[string]string{"message": "successfully processed"})
}
