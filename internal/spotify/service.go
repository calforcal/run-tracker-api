package spotify

import (
	"encoding/base64"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"run-tracker-api/internal/config"
	"strings"
	"time"

	"go.uber.org/zap"
)

type (
	SpotifyService struct {
		client *http.Client
		cfg    *config.Config
		logger *zap.Logger
	}

	TokenResponse struct {
		TokenType    string `json:"token_type"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
)

func New(cfg *config.Config, logger *zap.Logger) *SpotifyService {
	return &SpotifyService{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		cfg:    cfg,
		logger: logger,
	}
}

func (s *SpotifyService) ExchangeCodeForToken(code string, redirectURI string, grantType string) (TokenResponse, error) {
	clientID := s.cfg.SpotifyClientID
	clientSecret := s.cfg.SpotifyClientSecret

	encodedString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))

	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("redirect_uri", redirectURI)
	formData.Set("grant_type", grantType)

	encodedData := formData.Encode()
	reqBody := strings.NewReader(encodedData)

	url := "https://accounts.spotify.com/api/token"

	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return TokenResponse{}, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedString))

	resp, err := s.client.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}

	defer resp.Body.Close()

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return TokenResponse{}, nil
	}

	return tokenResponse, nil
}

func (s *SpotifyService) GetCurrentUser(accessToken string) (SpotifyUser, error) {
	url := "https://api.spotify.com/v1/me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return SpotifyUser{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println("response error", err)
		return SpotifyUser{}, err
	}

	defer resp.Body.Close()

	var spotifyUser SpotifyUser
	if err := json.NewDecoder(resp.Body).Decode(&spotifyUser); err != nil {
		return SpotifyUser{}, err
	}

	return spotifyUser, nil
}

func (s *SpotifyService) GetListeningHistory(accessToken string, after int64) (ListeningHistory, error) {
	baseURL := "https://api.spotify.com/v1/me/player/recently-played"

	// Build query parameters
	params := url.Values{}
	if after < 1 {
		params.Set("after", fmt.Sprintf("%d", after))
	}

	// Add limit parameter (Spotify API default is 20, max is 50)
	params.Set("limit", "50")

	fullURL := baseURL
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return ListeningHistory{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println("response error", err)
		return ListeningHistory{}, err
	}

	defer resp.Body.Close()

	// Read the response body once
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.Error(err))
		return ListeningHistory{}, err
	}

	var listeningHistory ListeningHistory
	if err := json.Unmarshal(bodyBytes, &listeningHistory); err != nil {
		s.logger.Error("Failed to decode response", zap.Error(err))
		return ListeningHistory{}, err
	}

	return listeningHistory, nil
}

func (s *SpotifyService) RefreshToken(refreshToken string) (TokenResponse, error) {
	clientID := s.cfg.SpotifyClientID
	clientSecret := s.cfg.SpotifyClientSecret

	encodedString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))

	baseUrl := "https://accounts.spotify.com/api/token"

	formData := url.Values{}
	formData.Set("refresh_token", refreshToken)
	formData.Set("client_id", clientID)
	formData.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("req is bad ", err)
		return TokenResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedString))

	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println("failed requesttt ", err)
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return TokenResponse{}, fmt.Errorf("spotify returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return TokenResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if tokenResponse.AccessToken == "" {
		return TokenResponse{}, fmt.Errorf("spotify returned empty access token")
	}

	return tokenResponse, nil
}
