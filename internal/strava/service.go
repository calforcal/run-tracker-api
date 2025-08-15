package strava

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	req.Header.Set("Authorization", "Bearer "+accessToken)
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

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return TokenResponse{}, err
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
