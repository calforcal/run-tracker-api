// Command main is the composition root: it loads config, builds every
// adapter, wires them into the application services, builds the HTTP
// handlers on top of those services, registers routes, and starts the
// server. No business logic lives here.
package main

import (
	"log"

	"run-tracker-api/api"
	"run-tracker-api/api/handlers/athlete"
	"run-tracker-api/api/handlers/auth"
	"run-tracker-api/api/handlers/home"
	"run-tracker-api/api/handlers/middleware"
	"run-tracker-api/api/handlers/user"
	"run-tracker-api/api/handlers/webhooks"
	"run-tracker-api/internal/adapters/postgres"
	spotifyadapter "run-tracker-api/internal/adapters/spotify"
	stravaadapter "run-tracker-api/internal/adapters/strava"
	"run-tracker-api/internal/config"
	activityservice "run-tracker-api/internal/services/activity"
	authservice "run-tracker-api/internal/services/auth"
	listeningservice "run-tracker-api/internal/services/listening"
	userservice "run-tracker-api/internal/services/user"
	webhookservice "run-tracker-api/internal/services/webhook"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	em "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	e := echo.New()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	logger, _ := zap.NewProduction()

	cfg := config.New()

	db := postgres.Connect(cfg, logger)

	userRepo := postgres.NewUserRepository(db)
	webhookRepo := postgres.NewWebhookRepository(db)
	listeningRepo := postgres.NewListeningHistoryRepository(db)

	stravaProvider := stravaadapter.NewProvider(cfg.StravaClientID, cfg.StravaClientSecret, nil)
	spotifyProvider := spotifyadapter.NewProvider(cfg.SpotifyClientID, cfg.SpotifyClientSecret, nil)

	authSvc := authservice.New(cfg.JwtSecret)
	userSvc := userservice.New(userRepo, stravaProvider, spotifyProvider)
	activitySvc := activityservice.New(stravaProvider)
	listeningSvc := listeningservice.New(spotifyProvider)
	webhookSvc := webhookservice.New(stravaProvider, spotifyProvider, userRepo, webhookRepo, listeningRepo, cfg.WebhookCallbackURL, cfg.WebhookToken)

	authMiddleware := middleware.NewAuthMiddleware(authSvc)
	adminMiddleware := middleware.NewAdminMiddleware(cfg.WebhookAdminToken)

	homeHandler := home.New()
	athleteHandler := athlete.New(userSvc, activitySvc)
	authHandler := auth.New(cfg, userSvc, authSvc, logger)
	userHandler := user.New(userSvc, listeningSvc, logger)
	webhookHandler := webhooks.New(logger, webhookSvc)

	e.Use(em.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			err := next(c)

			logger.Info("response",
				zap.Int("status", res.Status),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Error(err),
			)

			return err
		}
	})
	e.Use(em.CORSWithConfig(em.CORSConfig{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	}))

	api.RegisterRoutes(e, api.Handlers{
		Home:            homeHandler,
		Athlete:         athleteHandler,
		Auth:            authHandler,
		User:            userHandler,
		Webhook:         webhookHandler,
		AuthMiddleware:  authMiddleware,
		AdminMiddleware: adminMiddleware,
	})

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
