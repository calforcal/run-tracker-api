package config

import "os"

type Config struct {
	StravaAccessToken  string
	StravaClientID     string
	StravaClientSecret string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	MigrationsDir      string
	JwtSecret          string
}

func New() *Config {
	return &Config{
		StravaAccessToken:  os.Getenv("STRAVA_ACCESS_TOKEN"),
		StravaClientID:     os.Getenv("STRAVA_CLIENT_ID"),
		StravaClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		MigrationsDir:      os.Getenv("GOOSE_MIGRATION_DIR"),
		JwtSecret:          os.Getenv("JWT_SECRET"),
	}
}
