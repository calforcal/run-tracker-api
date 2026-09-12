package ports

import (
	"context"

	"run-tracker-api/internal/domain"
)

// StravaProvider is the driven port for the Strava external API.
type StravaProvider interface {
	ExchangeCodeForToken(ctx context.Context, code string) (domain.StravaCredentials, domain.Athlete, error)
	RefreshToken(ctx context.Context, refreshToken string) (domain.StravaCredentials, error)
	GetAthlete(ctx context.Context, accessToken string) (domain.Athlete, error)
	GetAthleteActivities(ctx context.Context, accessToken string) ([]domain.Activity, error)
	GetDetailedActivity(ctx context.Context, accessToken, activityID string) (domain.DetailedActivity, error)
	GetActivityStream(ctx context.Context, accessToken, activityID string) ([]domain.ActivityStream, error)
	CreateWebhookSubscription(ctx context.Context, callbackURL, verifyToken string) (domain.WebhookSubscription, error)
	ListWebhookSubscriptions(ctx context.Context) ([]domain.WebhookSubscription, error)
	DeleteWebhookSubscription(ctx context.Context, stravaSubscriptionID int) error
}

// SpotifyProvider is the driven port for the Spotify external API.
type SpotifyProvider interface {
	ExchangeCodeForToken(ctx context.Context, code, redirectURI string) (domain.SpotifyCredentials, error)
	RefreshToken(ctx context.Context, refreshToken string) (domain.SpotifyCredentials, error)
	GetCurrentUserID(ctx context.Context, accessToken string) (string, error)
	GetListeningHistory(ctx context.Context, accessToken string, after int64) ([]domain.ListeningHistoryItem, error)
}
