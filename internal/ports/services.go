package ports

import (
	"context"

	"run-tracker-api/internal/domain"
)

// AuthService issues and validates session tokens.
type AuthService interface {
	IssueToken(ctx context.Context, user domain.User) (string, error)
	ParseToken(ctx context.Context, token string) (domain.Claims, error)
}

// UserService owns account creation/linking and credential lifecycle.
type UserService interface {
	ExchangeStravaCode(ctx context.Context, code string) (domain.User, error)
	LoginWithStrava(ctx context.Context, code string) (domain.User, error)
	LinkSpotifyAccount(ctx context.Context, uuid, code, redirectURI string) (domain.User, error)
	GetByUUID(ctx context.Context, uuid string) (domain.User, error)
	EnsureValidStravaToken(ctx context.Context, user domain.User) (domain.User, error)
	EnsureValidSpotifyToken(ctx context.Context, user domain.User) (domain.User, error)
}

// ActivityService exposes a user's Strava athlete/activity data.
type ActivityService interface {
	GetAthlete(ctx context.Context, accessToken string) (domain.Athlete, error)
	GetActivities(ctx context.Context, accessToken string) ([]domain.Activity, error)
	GetActivity(ctx context.Context, accessToken, activityID string) (domain.DetailedActivity, error)
	GetActivityStream(ctx context.Context, accessToken, activityID string) ([]domain.ActivityStream, error)
}

// ListeningService exposes a user's Spotify listening history.
type ListeningService interface {
	GetListeningHistory(ctx context.Context, accessToken string, after int64) ([]domain.ListeningHistoryItem, error)
}

// WebhookService manages the Strava push-subscription lifecycle and
// correlates incoming activity events with Spotify listening history.
type WebhookService interface {
	CreateSubscription(ctx context.Context) (domain.WebhookSubscription, error)
	ListStravaSubscriptions(ctx context.Context) ([]domain.WebhookSubscription, error)
	DeleteSubscription(ctx context.Context) error
	VerifyCallback(ctx context.Context, hubMode, hubChallenge, hubVerifyToken string) (string, error)
	ProcessEvent(ctx context.Context, event domain.WebhookEvent) error
}
