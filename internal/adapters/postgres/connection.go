// Package postgres implements the repository ports against Postgres via
// database/sql + lib/pq.
package postgres

import (
	"database/sql"
	"fmt"

	"run-tracker-api/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"go.uber.org/zap"
)

// Connect opens a Postgres connection and runs pending goose migrations.
func Connect(cfg *config.Config, logger *zap.Logger) *sql.DB {
	logger.Info("Connecting to database...")
	connStr := cfg.DatabaseURL
	if connStr == "" {
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal("error starting database", zap.Error(err))
	}
	logger.Info("Pinging database...")
	if err := db.Ping(); err != nil {
		logger.Fatal("error pinging database", zap.Error(err))
	}
	logger.Info("Successfully connected to database")

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

	return db
}
