package webhooks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/spotify"
	"run-tracker-api/internal/storage"
	"run-tracker-api/internal/users"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	WebhookHandler struct {
		cfg            *config.Config
		client         *http.Client
		userService    *users.UserService
		spotifyService *spotify.SpotifyService
		storage        *storage.Storage
		logger         *zap.Logger
	}

	WebhookCreateResponse struct {
		ID int `json:"id"`
	}

	WebhookVerificationRequest struct {
		HubMode        string `query:"hub.mode"`
		HubChallenge   string `query:"hub.challenge"`
		HubVerifyToken string `query:"hub.verify_token"`
	}
)

func New(cfg *config.Config, userService *users.UserService, spotifyService *spotify.SpotifyService, storage *storage.Storage, logger *zap.Logger) WebhookHandler {
	return WebhookHandler{
		cfg: cfg,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		userService:    userService,
		spotifyService: spotifyService,
		storage:        storage,
		logger:         logger,
	}
}

func (h *WebhookHandler) CreateWebhook(c echo.Context) error {
	clientID := h.cfg.StravaClientID
	clientSecret := h.cfg.StravaClientSecret

	baseUrl := "https://www.strava.com/api/v3/push_subscriptions"
	callbackURL := "https://scotty-unglozed-nonvisibly.ngrok-free.dev/api/webhooks/strava"

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	params.Set("callback_url", callbackURL)
	params.Set("verify_token", h.cfg.WebhookToken)

	reqUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		h.logger.Info("error creating request %w", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error creating request"})
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := h.client.Do(req)
	if err != nil {
		h.logger.Info("error creating webhook with strava: %w", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "problem creating webhook upstream"})
	}

	defer resp.Body.Close()

	// Read the body once
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Info("error reading response body", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "problem reading response"})
	}

	// Print/log the response for debugging
	h.logger.Info("Strava API Response",
		zap.Int("status", resp.StatusCode),
		zap.String("body", string(bodyBytes)))

	// Check if status code indicates an error
	if resp.StatusCode > 399 {
		h.logger.Error("Strava API returned error",
			zap.Int("status", resp.StatusCode),
			zap.String("response", string(bodyBytes)))
		return c.JSON(resp.StatusCode, echo.Map{
			"error":   "strava API error",
			"details": string(bodyBytes),
			"status":  resp.StatusCode,
		})
	}

	var webhookResponse WebhookCreateResponse
	if err := json.Unmarshal(bodyBytes, &webhookResponse); err != nil {
		h.logger.Info("error decoding response: %w", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":        "problem decoding response",
			"raw_response": string(bodyBytes),
		})
	}

	fmt.Println(webhookResponse.ID)
	webhookSubscription, err := h.storage.CreateWebhookSubscription(webhookResponse.ID, callbackURL)
	if err != nil || webhookSubscription.ID == 0 {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("problem creating webhook subscription in database: %v", err)})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "webhook subscription created successfully"})
}

func (h *WebhookHandler) GetWebhook(c echo.Context) error {
	clientID := h.cfg.StravaClientID
	clientSecret := h.cfg.StravaClientSecret
	baseUrl := "https://www.strava.com/api/v3/push_subscriptions"

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)

	reqUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error creating request"})
	}

	resp, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "problem creating webhook upstream"})
	}

	defer resp.Body.Close()

	var webhookResponse []WebhookCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhookResponse); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("problem decoding response: %v", err)})
	}

	return c.JSON(http.StatusOK, webhookResponse)
}

func (h *WebhookHandler) DeleteWebhook(c echo.Context) error {
	clientID := h.cfg.StravaClientID
	clientSecret := h.cfg.StravaClientSecret

	webhookSubscription, err := h.storage.GetWebhookSubscription()
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": fmt.Sprintf("no webhook found to delete: %d", webhookSubscription.StravaID)})
	}
	baseUrl := fmt.Sprintf("https://www.strava.com/api/v3/push_subscriptions/%d", webhookSubscription.StravaID)

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)

	reqUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("DELETE", reqUrl, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error creating request"})
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "problem creating webhook upstream"})
	}

	defer resp.Body.Close()

	// Handle 204 No Content response
	if resp.StatusCode == http.StatusNoContent {
		err := h.storage.DeleteWebhook(webhookSubscription.StravaID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "problem deleting from database"})
		}
		return c.JSON(http.StatusNoContent, echo.Map{"message": "deleted successfully found"})
	}

	// Check for error status codes
	if resp.StatusCode > 399 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		h.logger.Error("API returned error",
			zap.Int("status", resp.StatusCode),
			zap.String("response", string(bodyBytes)))
		return c.JSON(resp.StatusCode, echo.Map{
			"error":   "API error",
			"details": string(bodyBytes),
			"status":  resp.StatusCode,
		})
	}

	var webhookResponse []WebhookCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhookResponse); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("problem decoding response: %v", err)})
	}

	return c.JSON(http.StatusOK, webhookResponse)
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
