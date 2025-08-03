package middleware

import (
	"run-tracker-api/internal/config"

	"go.uber.org/zap"
)

type (
	Middleware struct {
		cfg    *config.Config
		logger *zap.Logger
	}
)

func New(cfg *config.Config, logger *zap.Logger) *Middleware {
	return &Middleware{
		cfg:    cfg,
		logger: logger,
	}
}
