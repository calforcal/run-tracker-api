package storage

import (
	"database/sql"
	"fmt"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/spotify"
	"run-tracker-api/internal/strava"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"go.uber.org/zap"
)

type (
	Storage struct {
		cfg    *config.Config
		db     *sql.DB
		logger *zap.Logger
	}
)

func New(cfg *config.Config, logger *zap.Logger) *Storage {
	logger.Info("Connecting to database...")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal("error starting database", zap.Error(err))
	}
	logger.Info("Pinging database...")
	err = db.Ping()
	if err != nil {
		logger.Fatal("error pinging database", zap.Error(err))
	}
	logger.Info("Successfully connected to database")

	// Run goose migrations
	logger.Info("Running goose migrations...")
	goose.SetDialect("postgres")
	migrationsDir := cfg.MigrationsDir
	if migrationsDir == "" {
		logger.Fatal("Migrations directory is not set")
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		logger.Fatal("error running migrations", zap.Error(err))
	}
	logger.Info("Successfully ran goose migrations")

	return &Storage{cfg: cfg, db: db, logger: logger}
}

func (s *Storage) SaveUser(token *strava.TokenResponse) (User, error) {
	query :=
		`
		INSERT INTO users 
		(name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (strava_id) 
		DO UPDATE SET
				name = EXCLUDED.name,
				username = EXCLUDED.username,
				strava_access_token = EXCLUDED.strava_access_token,
				strava_refresh_token = EXCLUDED.strava_refresh_token,
				strava_expires_at = EXCLUDED.strava_expires_at,
				updated_at = NOW()
		RETURNING id, uuid, name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at, spotify_id, spotify_access_token, spotify_expires_at, spotify_refresh_token, created_at, updated_at;
	`

	var user User
	err := s.db.QueryRow(
		query,
		token.Athlete.Firstname+token.Athlete.Lastname,
		token.Athlete.Username,
		token.Athlete.ID,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresAt,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.StravaID,
		&user.StravaAccessToken,
		&user.StravaRefreshToken,
		&user.StravaExpiresAt,
		&user.SpotifyID,
		&user.SpotifyAccessToken,
		&user.SpotifyExpiresAt,
		&user.SpotifyRefreshToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return User{}, fmt.Errorf("error saving user: %w", err)
	}

	return user, nil
}

func (s *Storage) GetUserByStravaID(stravaId int64) (User, error) {
	user := User{}
	query := `SELECT id, uuid, name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at, created_at, updated_at FROM users WHERE strava_id = $1`
	if err := s.db.QueryRow(query, stravaId).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.StravaID,
		&user.StravaAccessToken,
		&user.StravaRefreshToken,
		&user.StravaExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Storage) UpdateUserFromToken(token *strava.TokenResponse) (User, error) {
	query := `
		UPDATE users SET
			name = $1,
			username = $2,
			strava_access_token = $3,
			strava_refresh_token = $4,
			strava_expires_at = $5,
			updated_at = NOW()
		WHERE strava_id = $6
		RETURNING id, uuid, name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at, created_at, updated_at;
	`

	var user User
	err := s.db.QueryRow(
		query,
		token.Athlete.Firstname+" "+token.Athlete.Lastname,
		token.Athlete.Username,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresAt,
		token.Athlete.ID,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.StravaID,
		&user.StravaAccessToken,
		&user.StravaRefreshToken,
		&user.StravaExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return User{}, fmt.Errorf("error saving user: %w", err)
	}

	return user, nil
}

func (s *Storage) GetUserByUUID(uuid string) (User, error) {
	user := User{}
	query := `SELECT id, uuid, name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at, spotify_id, spotify_access_token, spotify_expires_at, spotify_refresh_token, created_at, updated_at FROM users WHERE uuid = $1`
	if err := s.db.QueryRow(query, uuid).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.StravaID,
		&user.StravaAccessToken,
		&user.StravaRefreshToken,
		&user.StravaExpiresAt,
		&user.SpotifyID,
		&user.SpotifyAccessToken,
		&user.SpotifyExpiresAt,
		&user.SpotifyRefreshToken,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Storage) GetUserBySpotifyID(spotifyID string) (User, error) {
	var user User
	query := `SELECT id, uuid, name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at, spotify_id, spotify_access_token, spotify_expires_at, spotify_refresh_token, created_at, updated_at FROM users WHERE spotify_id = $1`
	if err := s.db.QueryRow(query, spotifyID).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.StravaID,
		&user.StravaAccessToken,
		&user.StravaRefreshToken,
		&user.StravaExpiresAt,
		&user.SpotifyID,
		&user.SpotifyAccessToken,
		&user.SpotifyExpiresAt,
		&user.SpotifyRefreshToken,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Storage) SaveSpotifyUser(tokenResponse spotify.TokenResponse, spotifyID string) (User, error) {
	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn))

	query := `
		INSERT INTO users
		(spotify_id, spotify_access_token, spotify_refresh_token, spotify_expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at, spotify_id, spotify_access_token, spotify_refresh_token, spotify_expires_at, created_at, updated_at;
		`

	var user User
	err := s.db.QueryRow(
		query,
		spotifyID,
		tokenResponse.AccessToken,
		tokenResponse.RefreshToken,
		expiresAt,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.StravaID,
		&user.StravaAccessToken,
		&user.StravaRefreshToken,
		&user.StravaExpiresAt,
		&user.SpotifyID,
		&user.SpotifyAccessToken,
		&user.SpotifyRefreshToken,
		&user.SpotifyExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return User{}, fmt.Errorf("error saving user: %w", err)
	}

	return user, nil
}

func (s *Storage) UpdateSpotifyUser(tokenResponse spotify.TokenResponse, spotifyID string) (User, error) {
	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn)).Unix()

	query :=
		`
			UPDATE users
			SET spotify_access_token = $1, spotify_expires_at = $2
			WHERE spotify_id = $3
			RETURNING id, uuid, name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at, spotify_id, spotify_access_token, spotify_refresh_token, spotify_expires_at, created_at, updated_at;
		`
	var user User
	err := s.db.QueryRow(
		query,
		tokenResponse.AccessToken,
		expiresAt,
		spotifyID,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.StravaID,
		&user.StravaAccessToken,
		&user.StravaRefreshToken,
		&user.StravaExpiresAt,
		&user.SpotifyID,
		&user.SpotifyAccessToken,
		&user.SpotifyRefreshToken,
		&user.SpotifyExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return User{}, fmt.Errorf("error saving user: %w", err)
	}

	return user, nil
}

func (s *Storage) AddSpotifyToStravaUser(tokenResponse spotify.TokenResponse, spotifyID string, uuid string) (User, error) {
	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second).Unix()

	query := `
		UPDATE users
		SET spotify_id = $1, spotify_access_token = $2, spotify_refresh_token = $3, spotify_expires_at = $4
		WHERE uuid = $5
		RETURNING id, uuid, name, username, strava_id, strava_access_token, strava_refresh_token, strava_expires_at, spotify_id, spotify_access_token, spotify_refresh_token, spotify_expires_at, created_at, updated_at;
		`

	var user User
	err := s.db.QueryRow(
		query,
		spotifyID,
		tokenResponse.AccessToken,
		tokenResponse.RefreshToken,
		expiresAt,
		uuid,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.StravaID,
		&user.StravaAccessToken,
		&user.StravaRefreshToken,
		&user.StravaExpiresAt,
		&user.SpotifyID,
		&user.SpotifyAccessToken,
		&user.SpotifyRefreshToken,
		&user.SpotifyExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return User{}, fmt.Errorf("error saving user: %w", err)
	}

	return user, nil
}

func (s *Storage) CreateWebhookSubscription(stravaID int, callbackURL string) (WebhookSubscription, error) {
	query := `INSERT INTO webhook_subscriptions (strava_id, callback_url) VALUES ($1, $2) RETURNING id, strava_id, callback_url`

	var webhookSubscription WebhookSubscription
	err := s.db.QueryRow(query, stravaID, callbackURL).Scan(
		&webhookSubscription.ID,
		&webhookSubscription.StravaID,
		&webhookSubscription.CallbackURL,
	)

	if err != nil {
		return WebhookSubscription{}, fmt.Errorf("error writing webhook subscription to database %w", err)
	}

	return webhookSubscription, nil
}

func (s *Storage) GetWebhookSubscription() (WebhookSubscription, error) {
	query := `SELECT id, strava_id, callback_url FROM webhook_subscriptions`

	var webhookSubscription WebhookSubscription
	err := s.db.QueryRow(query).Scan(
		&webhookSubscription.ID,
		&webhookSubscription.StravaID,
		&webhookSubscription.CallbackURL,
	)

	if err != nil {
		return WebhookSubscription{}, fmt.Errorf("error writing webhook subscription to database %w", err)
	}

	return webhookSubscription, nil
}

func (s *Storage) DeleteWebhook(stravaID int) error {
	query := `DELETE FROM webhook_subscriptions WHERE strava_id = $1`
	_, err := s.db.Exec(query, stravaID)
	if err != nil {
		return fmt.Errorf("error deleting webhook_subscription: %v", err)
	}
	return nil
}
