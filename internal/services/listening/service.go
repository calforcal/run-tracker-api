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
}

// New builds a ports.ListeningService.
func New(spotifyProvider ports.SpotifyProvider) ports.ListeningService {
	return &service{spotifyProvider: spotifyProvider}
}

var _ ports.ListeningService = (*service)(nil)

func (s *service) GetListeningHistory(ctx context.Context, accessToken string, after int64) ([]domain.ListeningHistoryItem, error) {
	return s.spotifyProvider.GetListeningHistory(ctx, accessToken, after)
}
