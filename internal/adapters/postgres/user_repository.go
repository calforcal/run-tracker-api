package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository builds a ports.UserRepository backed by Postgres.
func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &userRepository{db: db}
}

var _ ports.UserRepository = (*userRepository)(nil)

// returningColumns is the column list/order used consistently by every
// query in this file, paired with scanUser, so no query can silently drift
// out of sync with its Scan() call.
const returningColumns = `
	id, uuid, name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at,
	spotify_id, spotify_access_token, spotify_refresh_token, spotify_expires_at, created_at, updated_at
`

// scanner is satisfied by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...interface{}) error
}

func scanUser(s scanner) (domain.User, error) {
	var (
		u                                                  domain.User
		spotifyID, spotifyAccessToken, spotifyRefreshToken sql.NullString
		spotifyExpiresAt                                   sql.NullInt64
		stravaExpiresAt                                    int64
	)

	if err := s.Scan(
		&u.ID, &u.UUID, &u.Name, &u.Username,
		&u.Strava.AthleteID, &u.Strava.AccessToken, &u.Strava.RefreshToken, &stravaExpiresAt,
		&spotifyID, &spotifyAccessToken, &spotifyRefreshToken, &spotifyExpiresAt,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return domain.User{}, err
	}

	u.Strava.ExpiresAt = time.Unix(stravaExpiresAt, 0)

	if spotifyID.Valid && spotifyID.String != "" {
		u.Spotify = &domain.SpotifyCredentials{
			SpotifyID:    spotifyID.String,
			AccessToken:  spotifyAccessToken.String,
			RefreshToken: spotifyRefreshToken.String,
		}
		if spotifyExpiresAt.Valid {
			u.Spotify.ExpiresAt = time.Unix(spotifyExpiresAt.Int64, 0)
		}
	}

	return u, nil
}

func (r *userRepository) UpsertFromStrava(ctx context.Context, athlete domain.Athlete, creds domain.StravaCredentials) (domain.User, error) {
	query := fmt.Sprintf(`
		INSERT INTO users (name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (strava_id) DO UPDATE SET
			name = EXCLUDED.name,
			username = EXCLUDED.username,
			strava_access_token = EXCLUDED.strava_access_token,
			strava_refresh_token = EXCLUDED.strava_refresh_token,
			strava_expires_at = EXCLUDED.strava_expires_at,
			updated_at = NOW()
		RETURNING %s`, returningColumns)

	row := r.db.QueryRowContext(ctx, query,
		athlete.Firstname+athlete.Lastname,
		athlete.Username,
		creds.AthleteID,
		creds.AccessToken,
		creds.RefreshToken,
		creds.ExpiresAt.Unix(),
	)

	user, err := scanUser(row)
	if err != nil {
		return domain.User{}, fmt.Errorf("error upserting user from strava: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetByStravaID(ctx context.Context, stravaID int64) (domain.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users WHERE strava_id = $1`, returningColumns)
	user, err := scanUser(r.db.QueryRowContext(ctx, query, stravaID))
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetByUUID(ctx context.Context, uuid string) (domain.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users WHERE uuid = $1`, returningColumns)
	user, err := scanUser(r.db.QueryRowContext(ctx, query, uuid))
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetBySpotifyID(ctx context.Context, spotifyID string) (domain.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users WHERE spotify_id = $1`, returningColumns)
	user, err := scanUser(r.db.QueryRowContext(ctx, query, spotifyID))
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *userRepository) CreateSpotifyUser(ctx context.Context, spotifyID string, creds domain.SpotifyCredentials) (domain.User, error) {
	query := fmt.Sprintf(`
		INSERT INTO users (spotify_id, spotify_access_token, spotify_refresh_token, spotify_expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING %s`, returningColumns)

	row := r.db.QueryRowContext(ctx, query, spotifyID, creds.AccessToken, creds.RefreshToken, creds.ExpiresAt.Unix())

	user, err := scanUser(row)
	if err != nil {
		return domain.User{}, fmt.Errorf("error saving spotify user: %w", err)
	}
	return user, nil
}

func (r *userRepository) UpdateSpotifyCredentials(ctx context.Context, spotifyID string, creds domain.SpotifyCredentials) (domain.User, error) {
	query := fmt.Sprintf(`
		UPDATE users
		SET spotify_access_token = $1, spotify_refresh_token = $2, spotify_expires_at = $3, updated_at = NOW()
		WHERE spotify_id = $4
		RETURNING %s`, returningColumns)

	row := r.db.QueryRowContext(ctx, query, creds.AccessToken, creds.RefreshToken, creds.ExpiresAt.Unix(), spotifyID)

	user, err := scanUser(row)
	if err != nil {
		return domain.User{}, fmt.Errorf("error updating spotify credentials: %w", err)
	}
	return user, nil
}

func (r *userRepository) LinkSpotifyAccount(ctx context.Context, uuid string, spotifyID string, creds domain.SpotifyCredentials) (domain.User, error) {
	query := fmt.Sprintf(`
		UPDATE users
		SET spotify_id = $1, spotify_access_token = $2, spotify_refresh_token = $3, spotify_expires_at = $4, updated_at = NOW()
		WHERE uuid = $5
		RETURNING %s`, returningColumns)

	row := r.db.QueryRowContext(ctx, query, spotifyID, creds.AccessToken, creds.RefreshToken, creds.ExpiresAt.Unix(), uuid)

	user, err := scanUser(row)
	if err != nil {
		return domain.User{}, fmt.Errorf("error linking spotify account: %w", err)
	}
	return user, nil
}

func (r *userRepository) UpdateStravaCredentials(ctx context.Context, stravaID int64, creds domain.StravaCredentials) (domain.User, error) {
	query := fmt.Sprintf(`
		UPDATE users
		SET strava_access_token = $1, strava_refresh_token = $2, strava_expires_at = $3, updated_at = NOW()
		WHERE strava_id = $4
		RETURNING %s`, returningColumns)

	row := r.db.QueryRowContext(ctx, query, creds.AccessToken, creds.RefreshToken, creds.ExpiresAt.Unix(), stravaID)

	user, err := scanUser(row)
	if err != nil {
		return domain.User{}, fmt.Errorf("error updating strava credentials: %w", err)
	}
	return user, nil
}
