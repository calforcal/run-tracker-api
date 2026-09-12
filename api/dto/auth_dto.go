package dto

// ExchangeCodeRequest is the request body for /api/login,
// /api/strava/authorize-user, and /api/spotify/authorize-user.
type ExchangeCodeRequest struct {
	Code string `json:"code"`
}

// TokenResponse is the response body carrying an issued session token.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}
