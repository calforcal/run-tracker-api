package main

import (
	"log"
	"run-tracker-api/api/handlers/athlete"
	"run-tracker-api/api/handlers/auth"
	"run-tracker-api/api/handlers/home"
	authService "run-tracker-api/internal/auth"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/storage"
	"run-tracker-api/internal/strava"
	"run-tracker-api/internal/users"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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
	userService := users.New(config, logger, storage)
	authService := authService.New(config, logger)

	homeHandler := home.New()
	athleteHandler := athlete.New(config, stravaService, logger)
	authHandler := auth.New(config, stravaService, userService, authService, logger)

	api := e.Group("/api")

	api.GET("/home", homeHandler.Home)
	api.GET("/athlete", athleteHandler.GetAthlete)

	api.POST("/authorize-user", authHandler.AuthorizeStravaUser)

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
	e.Logger.Fatal(e.Start(":8000"))
}
