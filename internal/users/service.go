package users

import (
	"database/sql"
	"fmt"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/spotify"
	"run-tracker-api/internal/storage"
	"run-tracker-api/internal/strava"
	"strings"

	"go.uber.org/zap"
)

type (
	UserService struct {
		cfg            *config.Config
		logger         *zap.Logger
		storage        *storage.Storage
		spotifyService *spotify.SpotifyService
	}
)

func New(cfg *config.Config, logger *zap.Logger, storage *storage.Storage, spotifyService *spotify.SpotifyService) *UserService {
	return &UserService{cfg: cfg, logger: logger, storage: storage, spotifyService: spotifyService}
}

func (s *UserService) CreateOrUpdateUser(tokenResponse *strava.TokenResponse) (*storage.User, error) {
	user, err := s.storage.GetUserByStravaID(tokenResponse.Athlete.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			user, err = s.storage.SaveUser(tokenResponse)
			if err != nil {
				return nil, err
			}
			return &user, nil
		}
		return nil, err
	}
	user, err = s.storage.UpdateUserFromToken(tokenResponse)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) CreateOrUpdateSpotifyUser(tokenResponse *spotify.TokenResponse) (*storage.User, error) {
	spotifyUser, err := s.spotifyService.GetCurrentUser(tokenResponse.AccessToken)
	if err != nil {
		return &storage.User{}, err
	}

	spotifyID := spotifyUser.ID

	user, err := s.storage.GetUserBySpotifyID(spotifyID)
	if err != nil {
		if err == sql.ErrNoRows {
			user, err := s.storage.SaveSpotifyUser(*tokenResponse, spotifyID)
			if err != nil {
				return &storage.User{}, err
			}
			return &user, nil
		}
		return &storage.User{}, nil
	}

	user, err = s.storage.UpdateSpotifyUser(*tokenResponse, spotifyID)
	if err != nil {
		return &storage.User{}, nil
	}

	return &user, nil
}

func (s *UserService) GetUserByUUID(uuid string) (storage.User, error) {
	user, err := s.storage.GetUserByUUID(uuid)
	if err != nil {
		return storage.User{}, err
	}
	return user, nil
}

func (s *UserService) AddSpotifyToStravaUser(uuid string, tokenResponse *spotify.TokenResponse) (*storage.User, error) {
	tokenStr := strings.TrimPrefix(tokenResponse.AccessToken, "Bearer ")
	fmt.Println("ACC TOKIE", tokenResponse.AccessToken)
	tokenResponse.AccessToken = tokenStr
	fmt.Println("TRIMMED TOKEN", tokenStr)
	spotifyUser, err := s.spotifyService.GetCurrentUser(tokenResponse.AccessToken)
	if err != nil {
		fmt.Println("DAMN", err)
		return &storage.User{}, err
	}
	fmt.Println("SPOTIFY USER", spotifyUser)

	spotifyID := spotifyUser.ID

	user, err := s.storage.AddSpotifyToStravaUser(*tokenResponse, spotifyID, uuid)
	if err != nil {
		fmt.Println("DB ERR", err)
		return &storage.User{}, err
	}

	return &user, nil
}

func (s *UserService) UpdateSpotifyUser(user *storage.User, tokenResponse *spotify.TokenResponse) (*storage.User, error) {
	updatedUser, err := s.storage.UpdateSpotifyUser(*tokenResponse, *user.SpotifyID)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}
