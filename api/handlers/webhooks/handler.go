package webhooks

import (
	"context"
	"net/http"
	"time"

	"run-tracker-api/api/dto"
	"run-tracker-api/internal/ports"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// processTimeout bounds how long a background webhook processing job may
// run once detached from the originating request.
const processTimeout = 30 * time.Second

type WebhookHandler struct {
	logger         *zap.Logger
	webhookService ports.WebhookService
}

func New(logger *zap.Logger, webhookService ports.WebhookService) *WebhookHandler {
	return &WebhookHandler{logger: logger, webhookService: webhookService}
}

func (h *WebhookHandler) ProcessWebhooks(c echo.Context) error {
	var event dto.WebhookEventRequest
	if err := c.Bind(&event); err != nil {
		h.logger.Info("invalid webhook body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	domainEvent := event.ToDomain()
	detachedCtx := context.WithoutCancel(c.Request().Context())

	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.logger.Error("panic while processing webhook", zap.Any("panic", r))
			}
		}()

		ctx, cancel := context.WithTimeout(detachedCtx, processTimeout)
		defer cancel()

		if err := h.webhookService.ProcessEvent(ctx, domainEvent); err != nil {
			h.logger.Error("error processing webhook", zap.Error(err))
		}
	}()

	return c.JSON(http.StatusOK, nil)
}

func (h *WebhookHandler) CreateWebhook(c echo.Context) error {
	sub, err := h.webhookService.CreateSubscription(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "problem creating webhook: " + err.Error()})
	}
	return c.JSON(http.StatusCreated, dto.WebhookSubscriptionFromDomain(sub))
}

func (h *WebhookHandler) GetWebhook(c echo.Context) error {
	subs, err := h.webhookService.ListStravaSubscriptions(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error fetching webhook: " + err.Error()})
	}
	return c.JSON(http.StatusOK, dto.WebhookSubscriptionsFromDomain(subs))
}

func (h *WebhookHandler) DeleteWebhook(c echo.Context) error {
	if err := h.webhookService.DeleteSubscription(c.Request().Context()); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error deleting webhook"})
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (h *WebhookHandler) VerifyWebhookCallback(c echo.Context) error {
	var params dto.WebhookVerificationRequest
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid query parameters"})
	}

	challenge, err := h.webhookService.VerifyCallback(c.Request().Context(), params.HubMode, params.HubChallenge, params.HubVerifyToken)
	if err != nil {
		h.logger.Info("invalid webhook verification params", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid query parameters"})
	}

	return c.JSON(http.StatusOK, echo.Map{"hub.challenge": challenge})
}
