// Package api wires the HTTP route table onto an Echo instance. It is the
// only place route paths are declared.
package api

import (
	"run-tracker-api/api/handlers/athlete"
	"run-tracker-api/api/handlers/auth"
	"run-tracker-api/api/handlers/home"
	"run-tracker-api/api/handlers/middleware"
	"run-tracker-api/api/handlers/user"
	"run-tracker-api/api/handlers/webhooks"

	"github.com/labstack/echo/v4"
)

// Handlers bundles every HTTP handler and the auth middleware needed to
// register routes.
type Handlers struct {
	Home           *home.HomeHandler
	Athlete        *athlete.AthleteHandler
	Auth           *auth.AuthHandler
	User           *user.UserHandler
	Webhook        *webhooks.WebhookHandler
	AuthMiddleware *middleware.AuthMiddleware
}

// RegisterRoutes wires every route onto e.
func RegisterRoutes(e *echo.Echo, h Handlers) {
	e.GET("/", h.Home.Home)

	apiGroup := e.Group("/api")

	webhookGroup := apiGroup.Group("/webhooks")
	athleteGroup := apiGroup.Group("/athlete")
	userGroup := apiGroup.Group("/users")

	webhookGroup.GET("/strava/activity", h.Webhook.VerifyWebhookCallback)
	webhookGroup.POST("/strava/activity", h.Webhook.ProcessWebhooks)
	webhookGroup.POST("/strava", h.Webhook.CreateWebhook)
	webhookGroup.DELETE("/strava", h.Webhook.DeleteWebhook)
	webhookGroup.GET("/strava/view", h.Webhook.GetWebhook)

	userGroup.Use(h.AuthMiddleware.RunAuthMiddleware())
	userGroup.GET("/listening-history", h.User.GetListeningHistory)

	athleteGroup.Use(h.AuthMiddleware.RunAuthMiddleware())
	athleteGroup.GET("/activities", h.Athlete.GetAthleteActivities)
	athleteGroup.GET("/activities/:activity_id", h.Athlete.GetActivityByStravaId)
	athleteGroup.GET("/activities/:activity_id/stream", h.Athlete.GetActivityStream)
	athleteGroup.GET("", h.Athlete.GetAthlete)

	apiGroup.GET("/home", h.Home.Home)

	apiGroup.POST("/login", h.Auth.Login)
	apiGroup.POST("/strava/authorize-user", h.Auth.AuthorizeStravaUser)
	apiGroup.POST("/spotify/authorize-user", h.Auth.AuthorizeSpotifyUser)
}
