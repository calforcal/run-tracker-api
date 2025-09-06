package athlete

import (
	"database/sql"
	"fmt"
	"net/http"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/strava"
	"run-tracker-api/internal/users"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	AthleteHandler struct {
		config        *config.Config
		stravaService *strava.StravaService
		userService   *users.UserService
		logger        *zap.Logger
	}
)

func New(cfg *config.Config, stravaService *strava.StravaService, userService *users.UserService, logger *zap.Logger) *AthleteHandler {
	return &AthleteHandler{
		config:        cfg,
		stravaService: stravaService,
		logger:        logger,
		userService:   userService,
	}
}

func (h *AthleteHandler) GetAthlete(c echo.Context) error {
	uuid := c.Get("uuid").(string)
	user, err := h.userService.GetUserByUUID(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, "error: error getting user")
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	athlete, err := h.stravaService.GetAthlete(user.StravaAccessToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Println("USER STRUCT:", user)
	fmt.Println("SPOTIFY ID POINTER:", user.SpotifyID)
	if user.SpotifyID != nil {
		fmt.Println("SPOTIFY ID VALUE:", *user.SpotifyID)
	} else {
		fmt.Println("SPOTIFY ID IS NIL")
	}

	// Print other pointer fields to see their values
	if user.SpotifyAccessToken != nil {
		fmt.Println("SPOTIFY ACCESS TOKEN:", *user.SpotifyAccessToken)
	}
	if user.SpotifyRefreshToken != nil {
		fmt.Println("SPOTIFY REFRESH TOKEN:", *user.SpotifyRefreshToken)
	}
	if user.SpotifyExpiresAt != nil {
		fmt.Println("SPOTIFY EXPIRES AT:", *user.SpotifyExpiresAt)
	}
	if user.SpotifyID != nil && *user.SpotifyID != "" {
		spotifyConnected := true
		athlete.IsSpotifyConnected = &spotifyConnected
	} else {
		spotifyConnected := false
		athlete.IsSpotifyConnected = &spotifyConnected
	}

	return c.JSON(http.StatusOK, athlete)
}

func (h *AthleteHandler) GetAthleteActivities(c echo.Context) error {
	uuid := c.Get("uuid").(string)
	user, err := h.userService.GetUserByUUID(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, "error: error getting user")
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	activities, err := h.stravaService.GetAthleteActivities(user.StravaAccessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, activities)
}

func (h *AthleteHandler) GetActivityByStravaId(c echo.Context) error {
	uuid := c.Get("uuid").(string)
	activityId := c.Param("activity_id")

	user, err := h.userService.GetUserByUUID(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, "error: error getting user")
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	activity, err := h.stravaService.GetDetailedActivity(activityId, user.StravaAccessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, activity)
}
