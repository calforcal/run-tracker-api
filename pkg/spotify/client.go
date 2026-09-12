// Package spotify is a raw HTTP client for the Spotify Web API. It knows
// nothing about this application's domain model or business logic - it only
// speaks Spotify's wire format.
package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	accountsBaseURL = "https://accounts.spotify.com"
	apiBaseURL      = "https://api.spotify.com/v1"
)

// Client is a raw HTTP client for the Spotify Web API.
type Client struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
}

// New builds a Spotify API client. Pass nil for httpClient to use a
// sensible default.
func New(clientID, clientSecret string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{httpClient: httpClient, clientID: clientID, clientSecret: clientSecret}
}

func (c *Client) basicAuthHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.clientID, c.clientSecret)))
}

func (c *Client) ExchangeCodeForToken(ctx context.Context, code, redirectURI string) (TokenResponse, error) {
	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("redirect_uri", redirectURI)
	formData.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, accountsBaseURL+"/api/token", strings.NewReader(formData.Encode()))
	if err != nil {
		return TokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", c.basicAuthHeader())

	resp, err := c.httpClient.Do(req)
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

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error) {
	formData := url.Values{}
	formData.Set("refresh_token", refreshToken)
	formData.Set("client_id", c.clientID)
	formData.Set("grant_type", "refresh_token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, accountsBaseURL+"/api/token", strings.NewReader(formData.Encode()))
	if err != nil {
		return TokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", c.basicAuthHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
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

func (c *Client) GetCurrentUser(ctx context.Context, accessToken string) (SpotifyUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiBaseURL+"/me", nil)
	if err != nil {
		return SpotifyUser{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SpotifyUser{}, err
	}
	defer resp.Body.Close()

	var spotifyUser SpotifyUser
	if err := json.NewDecoder(resp.Body).Decode(&spotifyUser); err != nil {
		return SpotifyUser{}, err
	}

	return spotifyUser, nil
}

func (c *Client) GetListeningHistory(ctx context.Context, accessToken string, after int64) (ListeningHistory, error) {
	params := url.Values{}
	if after > 0 {
		params.Set("after", fmt.Sprintf("%d", after))
	}
	params.Set("limit", "50")

	reqURL := apiBaseURL + "/me/player/recently-played?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return ListeningHistory{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ListeningHistory{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ListeningHistory{}, err
	}

	var listeningHistory ListeningHistory
	if err := json.Unmarshal(bodyBytes, &listeningHistory); err != nil {
		return ListeningHistory{}, err
	}

	return listeningHistory, nil
}
