// Package spotify implements ports.SpotifyProvider by wrapping the raw
// pkg/spotify HTTP client and mapping its wire types to/from the domain
// model.
package spotify

import (
	"context"
	"net/http"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"
	spotifypkg "run-tracker-api/pkg/spotify"
)

type provider struct {
	client *spotifypkg.Client
}

// NewProvider builds a ports.SpotifyProvider backed by the real Spotify API.
func NewProvider(clientID, clientSecret string, httpClient *http.Client) ports.SpotifyProvider {
	return &provider{client: spotifypkg.New(clientID, clientSecret, httpClient)}
}

var _ ports.SpotifyProvider = (*provider)(nil)

func (p *provider) ExchangeCodeForToken(ctx context.Context, code, redirectURI string) (domain.SpotifyCredentials, error) {
	token, err := p.client.ExchangeCodeForToken(ctx, code, redirectURI)
	if err != nil {
		return domain.SpotifyCredentials{}, err
	}
	return toDomainCredentials(token), nil
}

func (p *provider) RefreshToken(ctx context.Context, refreshToken string) (domain.SpotifyCredentials, error) {
	token, err := p.client.RefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.SpotifyCredentials{}, err
	}
	return toDomainCredentials(token), nil
}

func (p *provider) GetCurrentUserID(ctx context.Context, accessToken string) (string, error) {
	user, err := p.client.GetCurrentUser(ctx, accessToken)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (p *provider) GetListeningHistory(ctx context.Context, accessToken string, after int64) ([]domain.ListeningHistoryItem, error) {
	history, err := p.client.GetListeningHistory(ctx, accessToken, after)
	if err != nil {
		return nil, err
	}
	return toDomainListeningHistory(history), nil
}
