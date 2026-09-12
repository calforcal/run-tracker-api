package domain

import "time"

// StravaCredentials holds a user's Strava OAuth token state.
type StravaCredentials struct {
	AthleteID    int64
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// SpotifyCredentials holds a user's Spotify OAuth token state.
type SpotifyCredentials struct {
	SpotifyID    string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// User is the core application user. Spotify is nil until the user links
// their Spotify account.
type User struct {
	ID        int
	UUID      string
	Name      string
	Username  string
	Strava    StravaCredentials
	Spotify   *SpotifyCredentials
	CreatedAt time.Time
	UpdatedAt time.Time
}
