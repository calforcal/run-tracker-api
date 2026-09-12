// Package user implements ports.UserService: account creation/linking and
// credential lifecycle. It depends only on domain + ports, never directly
// on adapters or pkg clients.
package user

import (
	"context"
	"errors"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"
	"run-tracker-api/internal/services/tokenrefresh"
)

type service struct {
	userRepo        ports.UserRepository
	stravaProvider  ports.StravaProvider
	spotifyProvider ports.SpotifyProvider
}

// New builds a ports.UserService.
func New(userRepo ports.UserRepository, stravaProvider ports.StravaProvider, spotifyProvider ports.SpotifyProvider) ports.UserService {
	return &service{userRepo: userRepo, stravaProvider: stravaProvider, spotifyProvider: spotifyProvider}
}

var _ ports.UserService = (*service)(nil)

func (s *service) ExchangeStravaCode(ctx context.Context, code string) (domain.User, error) {
	creds, athlete, err := s.stravaProvider.ExchangeCodeForToken(ctx, code)
	if err != nil {
		return domain.User{}, err
	}
	if athlete.ID <= 0 {
		return domain.User{}, errors.New("error authenticating user")
	}

	return s.userRepo.UpsertFromStrava(ctx, athlete, creds)
}

func (s *service) LoginWithStrava(ctx context.Context, code string) (domain.User, error) {
	user, err := s.ExchangeStravaCode(ctx, code)
	if err != nil {
		return domain.User{}, err
	}

	if user.Spotify != nil {
		user, err = s.EnsureValidSpotifyToken(ctx, user)
		if err != nil {
			return domain.User{}, err
		}
	}

	return user, nil
}

func (s *service) LinkSpotifyAccount(ctx context.Context, uuid, code, redirectURI string) (domain.User, error) {
	creds, err := s.spotifyProvider.ExchangeCodeForToken(ctx, code, redirectURI)
	if err != nil {
		return domain.User{}, err
	}

	spotifyID, err := s.spotifyProvider.GetCurrentUserID(ctx, creds.AccessToken)
	if err != nil {
		return domain.User{}, err
	}

	return s.userRepo.LinkSpotifyAccount(ctx, uuid, spotifyID, creds)
}

func (s *service) GetByUUID(ctx context.Context, uuid string) (domain.User, error) {
	return s.userRepo.GetByUUID(ctx, uuid)
}

func (s *service) EnsureValidStravaToken(ctx context.Context, user domain.User) (domain.User, error) {
	return tokenrefresh.EnsureValidStravaToken(ctx, user, s.userRepo, s.stravaProvider)
}

func (s *service) EnsureValidSpotifyToken(ctx context.Context, user domain.User) (domain.User, error) {
	return tokenrefresh.EnsureValidSpotifyToken(ctx, user, s.userRepo, s.spotifyProvider)
}
