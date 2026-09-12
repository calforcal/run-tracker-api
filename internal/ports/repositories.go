package ports

import (
	"context"

	"run-tracker-api/internal/domain"
)

// UserRepository persists and retrieves User aggregates.
type UserRepository interface {
	UpsertFromStrava(ctx context.Context, athlete domain.Athlete, creds domain.StravaCredentials) (domain.User, error)
	GetByStravaID(ctx context.Context, stravaID int64) (domain.User, error)
	GetByUUID(ctx context.Context, uuid string) (domain.User, error)
	GetBySpotifyID(ctx context.Context, spotifyID string) (domain.User, error)
	CreateSpotifyUser(ctx context.Context, spotifyID string, creds domain.SpotifyCredentials) (domain.User, error)
	UpdateSpotifyCredentials(ctx context.Context, spotifyID string, creds domain.SpotifyCredentials) (domain.User, error)
	LinkSpotifyAccount(ctx context.Context, uuid string, spotifyID string, creds domain.SpotifyCredentials) (domain.User, error)
	UpdateStravaCredentials(ctx context.Context, stravaID int64, creds domain.StravaCredentials) (domain.User, error)
}

// WebhookRepository persists the (single) Strava webhook subscription.
type WebhookRepository interface {
	Create(ctx context.Context, sub domain.WebhookSubscription) (domain.WebhookSubscription, error)
	Get(ctx context.Context) (domain.WebhookSubscription, error)
	Delete(ctx context.Context, stravaSubscriptionID int) error
}

// ListeningHistoryRepository persists a user's Spotify listening history,
// correlated to the Strava activity it was played during.
type ListeningHistoryRepository interface {
	SaveEntry(ctx context.Context, userID int, activityID int, item domain.ListeningHistoryItem) error
}
