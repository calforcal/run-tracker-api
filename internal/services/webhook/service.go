// Package webhook implements ports.WebhookService: managing the Strava
// push-subscription lifecycle and correlating incoming activity events with
// Spotify listening history.
package webhook

import (
	"context"
	"fmt"
	"strconv"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"
	"run-tracker-api/internal/services/tokenrefresh"
)

const (
	aspectTypeCreate   = "create"
	objectTypeActivity = "activity"
)

type service struct {
	stravaProvider  ports.StravaProvider
	spotifyProvider ports.SpotifyProvider
	userRepo        ports.UserRepository
	webhookRepo     ports.WebhookRepository
	listeningRepo   ports.ListeningHistoryRepository
	callbackURL     string
	verifyToken     string
}

// New builds a ports.WebhookService. callbackURL is the URL Strava will be
// told to POST events to; verifyToken is the shared secret used for
// Strava's subscription-validation handshake.
func New(
	stravaProvider ports.StravaProvider,
	spotifyProvider ports.SpotifyProvider,
	userRepo ports.UserRepository,
	webhookRepo ports.WebhookRepository,
	listeningRepo ports.ListeningHistoryRepository,
	callbackURL, verifyToken string,
) ports.WebhookService {
	return &service{
		stravaProvider:  stravaProvider,
		spotifyProvider: spotifyProvider,
		userRepo:        userRepo,
		webhookRepo:     webhookRepo,
		listeningRepo:   listeningRepo,
		callbackURL:     callbackURL,
		verifyToken:     verifyToken,
	}
}

var _ ports.WebhookService = (*service)(nil)

func (s *service) CreateSubscription(ctx context.Context) (domain.WebhookSubscription, error) {
	sub, err := s.stravaProvider.CreateWebhookSubscription(ctx, s.callbackURL, s.verifyToken)
	if err != nil {
		return domain.WebhookSubscription{}, fmt.Errorf("problem creating webhook: %w", err)
	}

	return s.webhookRepo.Create(ctx, sub)
}

func (s *service) ListStravaSubscriptions(ctx context.Context) ([]domain.WebhookSubscription, error) {
	return s.stravaProvider.ListWebhookSubscriptions(ctx)
}

func (s *service) DeleteSubscription(ctx context.Context) error {
	sub, err := s.webhookRepo.Get(ctx)
	if err != nil {
		return err
	}

	if err := s.stravaProvider.DeleteWebhookSubscription(ctx, sub.StravaID); err != nil {
		return err
	}

	return s.webhookRepo.Delete(ctx, sub.StravaID)
}

func (s *service) VerifyCallback(ctx context.Context, hubMode, hubChallenge, hubVerifyToken string) (string, error) {
	if hubVerifyToken != s.verifyToken || hubChallenge == "" {
		return "", fmt.Errorf("invalid verification parameters")
	}
	return hubChallenge, nil
}

func (s *service) ProcessEvent(ctx context.Context, event domain.WebhookEvent) error {
	if event.AspectType != aspectTypeCreate || event.ObjectType != objectTypeActivity {
		return nil
	}

	user, err := s.userRepo.GetByStravaID(ctx, event.OwnerID)
	if err != nil {
		return fmt.Errorf("error retrieving user from database: %w", err)
	}

	user, err = tokenrefresh.EnsureValidSpotifyToken(ctx, user, s.userRepo, s.spotifyProvider)
	if err != nil {
		return fmt.Errorf("error refreshing spotify token: %w", err)
	}

	user, err = tokenrefresh.EnsureValidStravaToken(ctx, user, s.userRepo, s.stravaProvider)
	if err != nil {
		return fmt.Errorf("error refreshing strava token: %w", err)
	}

	if user.Spotify == nil {
		// Nothing to correlate - the activity owner hasn't linked Spotify.
		return nil
	}

	activityID := strconv.Itoa(event.ObjectID)
	activity, err := s.stravaProvider.GetDetailedActivity(ctx, user.Strava.AccessToken, activityID)
	if err != nil {
		return fmt.Errorf("error getting activity from strava by id %s: %w", activityID, err)
	}

	afterMs := activity.StartDate.UnixMilli()

	history, err := s.spotifyProvider.GetListeningHistory(ctx, user.Spotify.AccessToken, afterMs)
	if err != nil {
		return fmt.Errorf("error getting user listening history: %w", err)
	}

	for _, item := range history {
		if err := s.listeningRepo.SaveEntry(ctx, user.ID, event.ObjectID, item); err != nil {
			return fmt.Errorf("error saving user listening history in database: %w", err)
		}
	}

	return nil
}
