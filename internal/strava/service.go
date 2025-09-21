package strava

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"run-tracker-api/internal/config"
	"time"

	"go.uber.org/zap"
)

type (
	StravaService struct {
		client *http.Client
		cfg    *config.Config
		logger *zap.Logger
	}

	TokenResponse struct {
		TokenType    string  `json:"token_type"`
		AccessToken  string  `json:"access_token"`
		RefreshToken string  `json:"refresh_token"`
		ExpiresAt    int     `json:"expires_at"`
		ExpiresIn    int     `json:"expires_in"`
		Athlete      Athlete `json:"athlete"`
	}

	RefreshTokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresAt    int    `json:"expires_at"`
		ExpiresIn    int    `json:"expires_in"`
	}

	RefreshRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RefreshToken string `json:"refresh_token"`
		GrantType    string `json:"grant_type"`
	}
)

func New(cfg *config.Config, logger *zap.Logger) *StravaService {
	return &StravaService{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		cfg:    cfg,
		logger: logger,
	}
}

func (s *StravaService) GetAthlete(accessToken string) (Athlete, error) {
	req, err := http.NewRequest("GET", "https://www.strava.com/api/v3/athlete", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	if err != nil {
		return Athlete{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return Athlete{}, err
	}
	defer resp.Body.Close()

	jsonDecoder := json.NewDecoder(resp.Body)
	var athlete Athlete
	err = jsonDecoder.Decode(&athlete)
	if err != nil {
		return Athlete{}, err
	}

	return athlete, nil
}

func (s *StravaService) GetAthleteActivities(accessToken string) ([]Activity, error) {
	req, err := http.NewRequest("GET", "https://www.strava.com/api/v3/athlete/activities", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if err != nil {
		return []Activity{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return []Activity{}, err
	}

	defer resp.Body.Close()

	jsonDecoder := json.NewDecoder(resp.Body)
	var activities []Activity
	err = jsonDecoder.Decode(&activities)
	if err != nil {
		return []Activity{}, err
	}

	return activities, nil
}

func (s *StravaService) GetDetailedActivity(activityId string, accessToken string) (DetailedActivity, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.strava.com/api/v3/activities/%s", activityId), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if err != nil {
		return DetailedActivity{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return DetailedActivity{}, err
	}

	defer resp.Body.Close()

	jsonDecoder := json.NewDecoder(resp.Body)
	var activity DetailedActivity
	err = jsonDecoder.Decode(&activity)
	if err != nil {
		return DetailedActivity{}, err
	}

	return activity, nil
}

func (s *StravaService) ExchangeCodeForToken(code string) (TokenResponse, error) {
	clientId := s.cfg.StravaClientID
	clientSecret := s.cfg.StravaClientSecret

	url := fmt.Sprintf(
		"https://www.strava.com/oauth/token?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code",
		clientId, clientSecret, code,
	)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return TokenResponse{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Log the failed response
		log.Printf("Strava API error - Status: %d, Body: %s", resp.StatusCode, string(body))
		return TokenResponse{}, fmt.Errorf("strava API returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		log.Printf("Failed to decode Strava response - Body: %s", string(body))
		return TokenResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return tokenResponse, nil
}

func (s *StravaService) RefreshToken(refreshToken string) (RefreshTokenResponse, error) {
	clientId := s.cfg.StravaClientID
	clientSecret := s.cfg.StravaClientSecret

	url := "https://www.strava.com/oauth/token"
	body := RefreshRequest{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return RefreshTokenResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return RefreshTokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	defer resp.Body.Close()

	var refreshResponse RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResponse); err != nil {
		return RefreshTokenResponse{}, err
	}

	return refreshResponse, nil
}
