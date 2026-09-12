package config

import (
	"os"
	"strings"
)

type Config struct {
	Port                string
	StravaAccessToken   string
	StravaClientID      string
	StravaClientSecret  string
	SpotifyClientID     string
	SpotifyClientSecret string
	DatabaseURL         string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
	MigrationsDir       string
	JwtSecret           string
	WebhookToken        string
	WebhookCallbackURL  string
	SpotifyRedirectURI  string
	CORSAllowedOrigins  []string
}

func New() *Config {
	return &Config{
		Port:                envOrDefault("PORT", "8000"),
		StravaAccessToken:   os.Getenv("STRAVA_ACCESS_TOKEN"),
		StravaClientID:      os.Getenv("STRAVA_CLIENT_ID"),
		StravaClientSecret:  os.Getenv("STRAVA_CLIENT_SECRET"),
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		DBHost:              os.Getenv("DB_HOST"),
		DBPort:              os.Getenv("DB_PORT"),
		DBUser:              os.Getenv("DB_USER"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBName:              os.Getenv("DB_NAME"),
		DBSSLMode:           envOrDefault("DB_SSLMODE", "disable"),
		MigrationsDir:       os.Getenv("GOOSE_MIGRATION_DIR"),
		JwtSecret:           os.Getenv("JWT_SECRET"),
		WebhookToken:        os.Getenv("WEBHOOK_TOKEN"),
		WebhookCallbackURL:  os.Getenv("WEBHOOK_CALLBACK_URL"),
		SpotifyRedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
		CORSAllowedOrigins:  splitAndTrim(envOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173")),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitAndTrim(csv string) []string {
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
