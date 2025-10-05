package config

import "os"

type Config struct {
	StravaAccessToken   string
	StravaClientID      string
	StravaClientSecret  string
	SpotifyClientID     string
	SpotifyClientSecret string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	MigrationsDir       string
	JwtSecret           string
	WebhookToken        string
}

func New() *Config {
	return &Config{
		StravaAccessToken:   os.Getenv("STRAVA_ACCESS_TOKEN"),
		StravaClientID:      os.Getenv("STRAVA_CLIENT_ID"),
		StravaClientSecret:  os.Getenv("STRAVA_CLIENT_SECRET"),
		SpotifyClientID:     os.Getenv("SPOTIfY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SCERET"),
		DBHost:              os.Getenv("DB_HOST"),
		DBPort:              os.Getenv("DB_PORT"),
		DBUser:              os.Getenv("DB_USER"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBName:              os.Getenv("DB_NAME"),
		MigrationsDir:       os.Getenv("GOOSE_MIGRATION_DIR"),
		JwtSecret:           os.Getenv("JWT_SECRET"),
		WebhookToken:        os.Getenv("WEBHOOK_TOKEN"),
	}
}
