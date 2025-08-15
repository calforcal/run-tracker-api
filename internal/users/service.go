package users

import (
	"database/sql"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/storage"
	"run-tracker-api/internal/strava"

	"go.uber.org/zap"
)

type (
	UserService struct {
		cfg     *config.Config
		logger  *zap.Logger
		storage *storage.Storage
	}
)

func New(cfg *config.Config, logger *zap.Logger, storage *storage.Storage) *UserService {
	return &UserService{cfg: cfg, logger: logger, storage: storage}
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

func (s *UserService) GetUserByUUID(uuid string) (storage.User, error) {
	user, err := s.storage.GetUserByUUID(uuid)
	if err != nil {
		return storage.User{}, err
	}
	return user, nil
}
