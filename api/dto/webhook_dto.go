package dto

import (
	"time"

	"run-tracker-api/internal/domain"
)

// WebhookVerificationRequest binds Strava's subscription-validation
// handshake query params.
type WebhookVerificationRequest struct {
	HubMode        string `query:"hub.mode"`
	HubChallenge   string `query:"hub.challenge"`
	HubVerifyToken string `query:"hub.verify_token"`
}

// WebhookEventRequest binds an inbound Strava webhook event body.
type WebhookEventRequest struct {
	AspectType     string `json:"aspect_type"`
	EventTime      int64  `json:"event_time"`
	ObjectID       int    `json:"object_id"`
	ObjectType     string `json:"object_type"`
	OwnerID        int64  `json:"owner_id"`
	SubscriptionID int    `json:"subscription_id"`
}

func (r WebhookEventRequest) ToDomain() domain.WebhookEvent {
	return domain.WebhookEvent{
		AspectType:     r.AspectType,
		EventTime:      time.Unix(r.EventTime, 0),
		ObjectID:       r.ObjectID,
		ObjectType:     r.ObjectType,
		OwnerID:        r.OwnerID,
		SubscriptionID: r.SubscriptionID,
	}
}

// WebhookSubscriptionResponse mirrors Strava's push_subscriptions response shape.
type WebhookSubscriptionResponse struct {
	ID int `json:"id"`
}

func WebhookSubscriptionFromDomain(s domain.WebhookSubscription) WebhookSubscriptionResponse {
	return WebhookSubscriptionResponse{ID: s.StravaID}
}

func WebhookSubscriptionsFromDomain(subs []domain.WebhookSubscription) []WebhookSubscriptionResponse {
	out := make([]WebhookSubscriptionResponse, len(subs))
	for i, s := range subs {
		out[i] = WebhookSubscriptionFromDomain(s)
	}
	return out
}
