package main

import (
	"log"
	"run-tracker-api/api/handlers/athlete"
	"run-tracker-api/api/handlers/auth"
	"run-tracker-api/api/handlers/home"
	"run-tracker-api/api/handlers/middleware"
	authService "run-tracker-api/internal/auth"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/storage"
	"run-tracker-api/internal/strava"
	"run-tracker-api/internal/users"

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
	userService := users.New(config, logger, storage)
	authService := authService.New(config, logger)

	authMiddleware := middleware.NewAuthMiddleware(config, authService)

	homeHandler := home.New()
	athleteHandler := athlete.New(config, stravaService, userService, logger)
	authHandler := auth.New(config, stravaService, userService, authService, logger)

	api := e.Group("/api")
	athlete := api.Group("/athlete")

	athlete.Use(authMiddleware.RunAuthMiddleware())
	athlete.GET("/activities", athleteHandler.GetAthleteActivities)
	athlete.GET("/activities/:activity_id", athleteHandler.GetActivityByStravaId)

	api.GET("/home", homeHandler.Home)
	athlete.GET("", athleteHandler.GetAthlete)

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
	e.Use(em.CORSWithConfig(em.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}, // Specify allowed HTTP methods
	}))
	e.Logger.Fatal(e.Start(":8000"))
}
