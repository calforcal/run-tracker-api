// Package strava implements ports.StravaProvider by wrapping the raw
// pkg/strava HTTP client and mapping its wire types to/from the domain
// model.
package strava

import (
	"context"
	"net/http"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"
	stravapkg "run-tracker-api/pkg/strava"
)

type provider struct {
	client *stravapkg.Client
}

// NewProvider builds a ports.StravaProvider backed by the real Strava API.
func NewProvider(clientID, clientSecret string, httpClient *http.Client) ports.StravaProvider {
	return &provider{client: stravapkg.New(clientID, clientSecret, httpClient)}
}

var _ ports.StravaProvider = (*provider)(nil)

func (p *provider) ExchangeCodeForToken(ctx context.Context, code string) (domain.StravaCredentials, domain.Athlete, error) {
	token, err := p.client.ExchangeCodeForToken(ctx, code)
	if err != nil {
		return domain.StravaCredentials{}, domain.Athlete{}, err
	}
	return toDomainCredentials(token), toDomainAthlete(token.Athlete), nil
}

func (p *provider) RefreshToken(ctx context.Context, refreshToken string) (domain.StravaCredentials, error) {
	refreshed, err := p.client.RefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.StravaCredentials{}, err
	}
	return toDomainRefreshCredentials(refreshed), nil
}

func (p *provider) GetAthlete(ctx context.Context, accessToken string) (domain.Athlete, error) {
	athlete, err := p.client.GetAthlete(ctx, accessToken)
	if err != nil {
		return domain.Athlete{}, err
	}
	return toDomainAthlete(athlete), nil
}

func (p *provider) GetAthleteActivities(ctx context.Context, accessToken string) ([]domain.Activity, error) {
	activities, err := p.client.GetAthleteActivities(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	return toDomainActivities(activities), nil
}

func (p *provider) GetDetailedActivity(ctx context.Context, accessToken, activityID string) (domain.DetailedActivity, error) {
	activity, err := p.client.GetDetailedActivity(ctx, accessToken, activityID)
	if err != nil {
		return domain.DetailedActivity{}, err
	}
	return toDomainDetailedActivity(activity), nil
}

func (p *provider) GetActivityStream(ctx context.Context, accessToken, activityID string) ([]domain.ActivityStream, error) {
	streams, err := p.client.GetActivityStream(ctx, accessToken, activityID)
	if err != nil {
		return nil, err
	}
	return toDomainActivityStreams(streams), nil
}

func (p *provider) CreateWebhookSubscription(ctx context.Context, callbackURL, verifyToken string) (domain.WebhookSubscription, error) {
	resp, err := p.client.CreateSubscription(ctx, callbackURL, verifyToken)
	if err != nil {
		return domain.WebhookSubscription{}, err
	}
	return domain.WebhookSubscription{StravaID: resp.ID, CallbackURL: callbackURL}, nil
}

func (p *provider) ListWebhookSubscriptions(ctx context.Context) ([]domain.WebhookSubscription, error) {
	subs, err := p.client.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.WebhookSubscription, len(subs))
	for i, s := range subs {
		out[i] = domain.WebhookSubscription{StravaID: s.ID}
	}
	return out, nil
}

func (p *provider) DeleteWebhookSubscription(ctx context.Context, stravaSubscriptionID int) error {
	return p.client.DeleteSubscription(ctx, stravaSubscriptionID)
}
