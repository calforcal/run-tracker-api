// Package listening implements ports.ListeningService, exposing a user's
// Spotify listening history.
package listening

import (
	"context"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"
)

type service struct {
	spotifyProvider ports.SpotifyProvider
	listeningRepo   ports.ListeningHistoryRepository
}

// New builds a ports.ListeningService.
func New(spotifyProvider ports.SpotifyProvider, listeningRepo ports.ListeningHistoryRepository) ports.ListeningService {
	return &service{spotifyProvider: spotifyProvider, listeningRepo: listeningRepo}
}

var _ ports.ListeningService = (*service)(nil)

func (s *service) GetListeningHistory(ctx context.Context, accessToken string, after int64) ([]domain.ListeningHistoryItem, error) {
	return s.spotifyProvider.GetListeningHistory(ctx, accessToken, after)
}

func (s *service) GetHistoryForActivity(ctx context.Context, userID, activityID int) ([]domain.ListeningHistoryItem, error) {
	return s.listeningRepo.GetForActivity(ctx, userID, activityID)
}
