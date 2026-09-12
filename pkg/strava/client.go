// Package strava is a raw HTTP client for the Strava v3 API. It knows
// nothing about this application's domain model or business logic - it only
// speaks Strava's wire format.
package strava

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL       = "https://www.strava.com/api/v3"
	oauthTokenURL = "https://www.strava.com/oauth/token"
)

// Client is a raw HTTP client for the Strava API.
type Client struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
}

// New builds a Strava API client. Pass nil for httpClient to use a
// sensible default.
func New(clientID, clientSecret string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{httpClient: httpClient, clientID: clientID, clientSecret: clientSecret}
}

func (c *Client) ExchangeCodeForToken(ctx context.Context, code string) (TokenResponse, error) {
	reqURL := fmt.Sprintf(
		"%s/oauth/token?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code",
		baseURL, c.clientID, c.clientSecret, code,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return TokenResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return TokenResponse{}, fmt.Errorf("strava API returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return TokenResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return tokenResponse, nil
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (RefreshTokenResponse, error) {
	body := RefreshRequest{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return RefreshTokenResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oauthTokenURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return RefreshTokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
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

func (c *Client) GetAthlete(ctx context.Context, accessToken string) (Athlete, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/athlete", nil)
	if err != nil {
		return Athlete{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Athlete{}, err
	}
	defer resp.Body.Close()

	var athlete Athlete
	if err := json.NewDecoder(resp.Body).Decode(&athlete); err != nil {
		return Athlete{}, err
	}

	return athlete, nil
}

func (c *Client) GetAthleteActivities(ctx context.Context, accessToken string) ([]Activity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/athlete/activities", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var activities []Activity
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, err
	}

	return activities, nil
}

func (c *Client) GetDetailedActivity(ctx context.Context, accessToken, activityID string) (DetailedActivity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/activities/%s", baseURL, activityID), nil)
	if err != nil {
		return DetailedActivity{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return DetailedActivity{}, err
	}
	defer resp.Body.Close()

	var activity DetailedActivity
	if err := json.NewDecoder(resp.Body).Decode(&activity); err != nil {
		return DetailedActivity{}, err
	}

	return activity, nil
}

func (c *Client) GetActivityStream(ctx context.Context, accessToken, activityID string) ([]ActivityStream, error) {
	keysParam := "time,latlng,altitude,heartrate,watts"
	reqURL := fmt.Sprintf("%s/activities/%s/streams?keys=%s", baseURL, activityID, keysParam)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("strava API returned status %d", resp.StatusCode)
	}

	var streams []ActivityStream
	if err := json.NewDecoder(resp.Body).Decode(&streams); err != nil {
		return nil, err
	}

	return streams, nil
}

func (c *Client) CreateSubscription(ctx context.Context, callbackURL, verifyToken string) (WebhookSubscriptionResponse, error) {
	params := url.Values{}
	params.Set("client_id", c.clientID)
	params.Set("client_secret", c.clientSecret)
	params.Set("callback_url", callbackURL)
	params.Set("verify_token", verifyToken)

	reqURL := fmt.Sprintf("%s/push_subscriptions?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return WebhookSubscriptionResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return WebhookSubscriptionResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return WebhookSubscriptionResponse{}, err
	}

	if resp.StatusCode > 399 {
		return WebhookSubscriptionResponse{}, fmt.Errorf("strava API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var subscription WebhookSubscriptionResponse
	if err := json.Unmarshal(bodyBytes, &subscription); err != nil {
		return WebhookSubscriptionResponse{}, err
	}

	return subscription, nil
}

func (c *Client) ListSubscriptions(ctx context.Context) ([]WebhookSubscriptionResponse, error) {
	params := url.Values{}
	params.Set("client_id", c.clientID)
	params.Set("client_secret", c.clientSecret)

	reqURL := fmt.Sprintf("%s/push_subscriptions?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var subscriptions []WebhookSubscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (c *Client) DeleteSubscription(ctx context.Context, stravaSubscriptionID int) error {
	params := url.Values{}
	params.Set("client_id", c.clientID)
	params.Set("client_secret", c.clientSecret)

	reqURL := fmt.Sprintf("%s/push_subscriptions/%d?%s", baseURL, stravaSubscriptionID, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if resp.StatusCode > 399 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("strava API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
