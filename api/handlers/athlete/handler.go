package athlete

import (
	"database/sql"
	"net/http"
	"strconv"

	"run-tracker-api/api/dto"
	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"

	"github.com/labstack/echo/v4"
)

type AthleteHandler struct {
	userService      ports.UserService
	activityService  ports.ActivityService
	listeningService ports.ListeningService
}

func New(userService ports.UserService, activityService ports.ActivityService, listeningService ports.ListeningService) *AthleteHandler {
	return &AthleteHandler{userService: userService, activityService: activityService, listeningService: listeningService}
}

func (h *AthleteHandler) GetAthlete(c echo.Context) error {
	ctx := c.Request().Context()
	uuid := c.Get("uuid").(string)

	user, err := h.userService.GetByUUID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "error getting user"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	athlete, err := h.activityService.GetAthlete(ctx, user.Strava.AccessToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dto.AthleteFromDomain(athlete, user.Spotify != nil))
}

func (h *AthleteHandler) GetAthleteActivities(c echo.Context) error {
	ctx := c.Request().Context()
	uuid := c.Get("uuid").(string)

	user, err := h.userService.GetByUUID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "error getting user"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	activities, err := h.activityService.GetActivities(ctx, user.Strava.AccessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ActivitiesFromDomain(activities))
}

func (h *AthleteHandler) GetActivityByStravaId(c echo.Context) error {
	ctx := c.Request().Context()
	uuid := c.Get("uuid").(string)
	activityId := c.Param("activity_id")

	user, err := h.userService.GetByUUID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "error getting user"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	activity, err := h.activityService.GetActivity(ctx, user.Strava.AccessToken, activityId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	songs := []domain.ListeningHistoryItem{}
	if activityIDInt, parseErr := strconv.Atoi(activityId); parseErr == nil {
		songs, err = h.listeningService.GetHistoryForActivity(ctx, user.ID, activityIDInt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, dto.DetailedActivityFromDomain(activity, songs))
}

func (h *AthleteHandler) GetActivityStream(c echo.Context) error {
	ctx := c.Request().Context()
	uuid := c.Get("uuid").(string)
	activityId := c.Param("activity_id")

	user, err := h.userService.GetByUUID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "error getting user"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	stream, err := h.activityService.GetActivityStream(ctx, user.Strava.AccessToken, activityId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ActivityStreamsFromDomain(stream))
}
