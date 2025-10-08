package webhooks

import (
	"fmt"
	"net/http"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/webhooks"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	WebhookHandler struct {
		cfg            *config.Config
		logger         *zap.Logger
		webhookService *webhooks.WebhookService
	}

	WebhookVerificationRequest struct {
		HubMode        string `query:"hub.mode"`
		HubChallenge   string `query:"hub.challenge"`
		HubVerifyToken string `query:"hub.verify_token"`
	}
)

const (
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"

	ACTIVITY = "activity"
	ATHLETE  = "athlete"
)

func New(cfg *config.Config, logger *zap.Logger, webhookService *webhooks.WebhookService) WebhookHandler {
	return WebhookHandler{
		cfg:            cfg,
		logger:         logger,
		webhookService: webhookService,
	}
}

func (h *WebhookHandler) ProcessWebhooks(c echo.Context) error {
	var event webhooks.WebhookEvent
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	if event.AspectType == CREATE {
		if event.ObjectType == ACTIVITY {
			err := h.webhookService.ProcessActivity(event)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error processing webhook"})
			}
		}
		return c.JSON(http.StatusOK, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *WebhookHandler) CreateWebhook(c echo.Context) error {
	webhookSubscription, err := h.webhookService.CreateWebhook()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("problem creating webhook: %v", err)})
	}
	return c.JSON(http.StatusCreated, webhookSubscription)
}

func (h *WebhookHandler) GetWebhook(c echo.Context) error {
	webhookResponse, err := h.webhookService.GetWebhook()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("error fetching webhook: %v", err)})
	}

	return c.JSON(http.StatusOK, webhookResponse)
}

func (h *WebhookHandler) DeleteWebhook(c echo.Context) error {
	err := h.webhookService.DeleteWebhook()
	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "error deleting webhook"})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *WebhookHandler) VerifyWebhookCallback(c echo.Context) error {
	var params WebhookVerificationRequest

	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid query parameters",
		})
	}

	if params.HubVerifyToken != h.cfg.WebhookToken || params.HubChallenge == "" {
		h.logger.Info(fmt.Sprintf("invalid params: token: %s, challenge: %s", params.HubVerifyToken, params.HubChallenge))
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid query parameters",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"hub.challenge": params.HubChallenge,
	})
}

func (h *WebhookHandler) ActivityUpload(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]string{"message": "successfully processed"})
}
