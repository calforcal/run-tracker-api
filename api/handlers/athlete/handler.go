package athlete

import (
	"net/http"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/strava"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	AthleteHandler struct {
		config        *config.Config
		stravaService *strava.StravaService
		logger        *zap.Logger
	}
)

func New(cfg *config.Config, stravaService *strava.StravaService, logger *zap.Logger) *AthleteHandler {
	return &AthleteHandler{
		config:        cfg,
		stravaService: stravaService,
		logger:        logger,
	}
}

func (h *AthleteHandler) GetAthlete(c echo.Context) error {
	athlete, err := h.stravaService.GetAthlete(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, athlete)
}
