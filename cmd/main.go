package main

import (
	"log"
	"run-tracker-api/api/handlers/athlete"
	"run-tracker-api/api/handlers/auth"
	"run-tracker-api/api/handlers/home"
	"run-tracker-api/api/handlers/middleware"
	"run-tracker-api/api/handlers/user"
	"run-tracker-api/api/handlers/webhooks"
	authService "run-tracker-api/internal/auth"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/spotify"
	"run-tracker-api/internal/storage"
	"run-tracker-api/internal/strava"
	"run-tracker-api/internal/users"
	whs "run-tracker-api/internal/webhooks"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	em "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger, _ := zap.NewProduction()

	config := config.New()

	storage := storage.New(config, logger)

	stravaService := strava.New(config, logger)
	spotifyService := spotify.New(config, logger)
	userService := users.New(config, logger, storage, spotifyService)
	authService := authService.New(config, logger)
	webhookService := whs.New(config, logger, spotifyService, storage, stravaService, userService)

	authMiddleware := middleware.NewAuthMiddleware(config, authService)

	homeHandler := home.New()
	athleteHandler := athlete.New(config, stravaService, userService, logger)
	authHandler := auth.New(config, stravaService, spotifyService, userService, authService, logger)
	userHandler := user.New(config, spotifyService, userService, logger)

	wh := webhooks.New(config, logger, webhookService)

	api := e.Group("/api")

	webhook := api.Group("/webhooks")
	athlete := api.Group("/athlete")
	user := api.Group("/users")

	webhook.GET("/strava", wh.VerifyWebhookCallback)
	webhook.POST("/strava", wh.CreateWebhook)
	webhook.DELETE("/strava", wh.DeleteWebhook)
	webhook.GET("/strava/view", wh.GetWebhook)

	user.Use(authMiddleware.RunAuthMiddleware())

	user.GET("/listening-history", userHandler.GetListeningHistory)

	athlete.Use(authMiddleware.RunAuthMiddleware())
	athlete.GET("/activities", athleteHandler.GetAthleteActivities)
	athlete.GET("/activities/:activity_id", athleteHandler.GetActivityByStravaId)

	api.GET("/home", homeHandler.Home)
	athlete.GET("", athleteHandler.GetAthlete)

	api.POST("/login", authHandler.Login)
	api.POST("/strava/authorize-user", authHandler.AuthorizeStravaUser)
	api.POST("/spotify/authorize-user", authHandler.AuthorizeSpotifyUser)

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
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}, // Specify allowed HTTP methods
	}))
	e.Logger.Fatal(e.Start(":8000"))
}
