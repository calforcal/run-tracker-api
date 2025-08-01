package main

import (
	"log"
	"run-tracker-api/api/handlers/athlete"
	"run-tracker-api/api/handlers/home"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/strava"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := config.New()
	stravaService := strava.New()

	homeHandler := home.New()
	athleteHandler := athlete.New(config, stravaService)

	api := e.Group("/api")

	api.GET("/home", homeHandler.Home)
	api.GET("/athlete", athleteHandler.GetAthlete)

	e.Logger.Fatal(e.Start(":8000"))
}
