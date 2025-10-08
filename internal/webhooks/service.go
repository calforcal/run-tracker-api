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
	"run-tracker-api/internal/strava"
	"run-tracker-api/internal/users"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type (
	WebhookService struct {
		cfg            *config.Config
		client         *http.Client
		logger         *zap.Logger
		spotifyService *spotify.SpotifyService
		storage        *storage.Storage
		stravaService  *strava.StravaService
		usersService   *users.UserService
	}

	WebhookResponse struct {
		ID int `json:"id"`
	}

	WebhookEvent struct {
		AspectType     string `json:"aspect_type"`
		EventTime      string `json:"event_time"`
		ObjectID       int    `json:"object_id"`
		ObjectType     string `json:"object_type"`
		OwnerID        int64  `json:"owner_id"`
		SubscriptionID string `json:"subscription_id"`
	}
)

func New(cfg *config.Config, logger *zap.Logger, spotifyService *spotify.SpotifyService, storage *storage.Storage, stravaService *strava.StravaService, usersService *users.UserService) *WebhookService {
	return &WebhookService{
		cfg: cfg,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		logger:         logger,
		spotifyService: spotifyService,
		storage:        storage,
		stravaService:  stravaService,
		usersService:   usersService,
	}
}

func (s *WebhookService) CreateWebhook() (WebhookResponse, error) {
	clientID := s.cfg.StravaClientID
	clientSecret := s.cfg.StravaClientSecret

	baseUrl := "https://www.strava.com/api/v3/push_subscriptions"
	callbackURL := "https://scotty-unglozed-nonvisibly.ngrok-free.dev/api/webhooks/strava"

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	params.Set("callback_url", callbackURL)
	params.Set("verify_token", s.cfg.WebhookToken)

	reqUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		s.logger.Info("error creating request %w", zap.Error(err))
		return WebhookResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Info("error creating webhook with strava: %w", zap.Error(err))
		return WebhookResponse{}, err
	}

	defer resp.Body.Close()

	// Read the body once
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Info("error reading response body", zap.Error(err))
		return WebhookResponse{}, err
	}

	// Print/log the response for debugging
	s.logger.Info("Strava API Response",
		zap.Int("status", resp.StatusCode),
		zap.String("body", string(bodyBytes)))

	// Check if status code indicates an error
	if resp.StatusCode > 399 {
		s.logger.Error("Strava API returned error",
			zap.Int("status", resp.StatusCode),
			zap.String("response", string(bodyBytes)))
		return WebhookResponse{}, err
	}

	var webhookResponse WebhookResponse
	if err := json.Unmarshal(bodyBytes, &webhookResponse); err != nil {
		s.logger.Info("error decoding response: %w", zap.Error(err))
		return WebhookResponse{}, err
	}

	fmt.Println(webhookResponse.ID)
	webhookSubscription, err := s.storage.CreateWebhookSubscription(webhookResponse.ID, callbackURL)
	if err != nil || webhookSubscription.ID == 0 {
		s.logger.Info(fmt.Sprintf("error creating webhook subscription in database: %v", err))
		return WebhookResponse{}, err
	}
	return webhookResponse, nil
}

func (s *WebhookService) GetWebhook() ([]WebhookResponse, error) {
	clientID := s.cfg.StravaClientID
	clientSecret := s.cfg.StravaClientSecret
	baseUrl := "https://www.strava.com/api/v3/push_subscriptions"

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)

	reqUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error creating request: %v", err))
		return []WebhookResponse{}, nil
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error with request to strava: %v", err))
		return []WebhookResponse{}, nil
	}

	defer resp.Body.Close()

	var webhookResponse []WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhookResponse); err != nil {
		s.logger.Info(fmt.Sprintf("error decoding request: %v", err))
		return []WebhookResponse{}, nil
	}

	return webhookResponse, nil
}

func (s *WebhookService) DeleteWebhook() error {
	clientID := s.cfg.StravaClientID
	clientSecret := s.cfg.StravaClientSecret

	webhookSubscription, err := s.storage.GetWebhookSubscription()
	if err != nil {
		s.logger.Info(fmt.Sprintf("error getting webhook subscription from database: %v", err))
		return err
	}
	baseUrl := fmt.Sprintf("https://www.strava.com/api/v3/push_subscriptions/%d", webhookSubscription.StravaID)

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)

	reqUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("DELETE", reqUrl, nil)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error creating get webhook request to strava: %v", err))
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error deleting webhook with strava: %v", err))
		return err
	}

	defer resp.Body.Close()

	// Handle 204 No Content response
	if resp.StatusCode == http.StatusNoContent {
		err := s.storage.DeleteWebhook(webhookSubscription.StravaID)
		if err != nil {
			s.logger.Info(fmt.Sprintf("error deleting webhook in database: %v", err))
			return err
		}
		return nil
	}

	// Check for error status codes
	if resp.StatusCode > 399 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		s.logger.Error("API returned error",
			zap.Int("status", resp.StatusCode),
			zap.String("response", string(bodyBytes)))
		return fmt.Errorf("error deleting webhook with strava")
	}

	return nil
}

func (s *WebhookService) ProcessActivity(event WebhookEvent) error {
	user, err := s.storage.GetUserByStravaID(event.OwnerID)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error retrieving user from database: %v", err))
		return err
	}

	tokenResponse, err := s.spotifyService.RefreshToken(*user.SpotifyRefreshToken)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error refreshing token: %v", err))
		return err
	}

	updatedUser, err := s.usersService.UpdateSpotifyUser(&user, &tokenResponse)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error updating user in database: %v", err))
		return err
	}

	// Going to put this in here with the calls.
	// Maybe not optimal but holding off on a cron job for now to not "over engineer"
	stringId := strconv.Itoa(event.ObjectID)
	activity, err := s.stravaService.GetDetailedActivity(stringId, *updatedUser.SpotifyAccessToken)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error getting activity from strava by id -> %s: %v", stringId, err))
		return err
	}

	t, err := time.Parse(time.RFC3339, activity.StartDate)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error parsing time to ms: %v", err))
		return err
	}

	// Milliseconds
	unixMs := t.UnixMilli()

	spotifyToken := updatedUser.SpotifyAccessToken

	listeningHistory, err := s.spotifyService.GetListeningHistory(*spotifyToken, unixMs)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error getting user listening history: %v", err))
		return err
	}

	err = s.storage.SaveListeningHistory(listeningHistory, activity.ID)
	if err != nil {
		s.logger.Info(fmt.Sprintf("error saving user listening history in database: %v", err))
		return err
	}

	return nil
}
