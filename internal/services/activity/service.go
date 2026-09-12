// Package activity implements ports.ActivityService, exposing a user's
// Strava athlete/activity data.
package activity

import (
	"context"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"
)

type service struct {
	stravaProvider ports.StravaProvider
}

// New builds a ports.ActivityService.
func New(stravaProvider ports.StravaProvider) ports.ActivityService {
	return &service{stravaProvider: stravaProvider}
}

var _ ports.ActivityService = (*service)(nil)

func (s *service) GetAthlete(ctx context.Context, accessToken string) (domain.Athlete, error) {
	return s.stravaProvider.GetAthlete(ctx, accessToken)
}

func (s *service) GetActivities(ctx context.Context, accessToken string) ([]domain.Activity, error) {
	return s.stravaProvider.GetAthleteActivities(ctx, accessToken)
}

func (s *service) GetActivity(ctx context.Context, accessToken, activityID string) (domain.DetailedActivity, error) {
	return s.stravaProvider.GetDetailedActivity(ctx, accessToken, activityID)
}

func (s *service) GetActivityStream(ctx context.Context, accessToken, activityID string) ([]domain.ActivityStream, error) {
	return s.stravaProvider.GetActivityStream(ctx, accessToken, activityID)
}
