package domain

import "time"

// WebhookSubscription is a registered Strava push-subscription.
type WebhookSubscription struct {
	ID          int
	StravaID    int
	CallbackURL string
}

// WebhookEvent is an inbound Strava webhook event.
type WebhookEvent struct {
	AspectType     string
	EventTime      time.Time
	ObjectID       int
	ObjectType     string
	OwnerID        int64
	SubscriptionID int
}
