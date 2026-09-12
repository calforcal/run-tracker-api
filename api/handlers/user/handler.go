package user

import (
	"net/http"

	"run-tracker-api/api/dto"
	"run-tracker-api/internal/ports"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService      ports.UserService
	listeningService ports.ListeningService
	logger           *zap.Logger
}

func New(userService ports.UserService, listeningService ports.ListeningService, logger *zap.Logger) *UserHandler {
	return &UserHandler{userService: userService, listeningService: listeningService, logger: logger}
}

func (h *UserHandler) GetListeningHistory(c echo.Context) error {
	ctx := c.Request().Context()
	uuid := c.Get("uuid").(string)

	var params dto.ListeningHistoryRequest
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request parameters"})
	}

	user, err := h.userService.GetByUUID(ctx, uuid)
	if err != nil {
		h.logger.Info("no user found for uuid", zap.String("uuid", uuid))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "error getting user"})
	}

	if user.Spotify == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user has not linked a spotify account"})
	}

	user, err = h.userService.EnsureValidSpotifyToken(ctx, user)
	if err != nil {
		h.logger.Error("failed to refresh spotify token", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to refresh token"})
	}

	history, err := h.listeningService.GetListeningHistory(ctx, user.Spotify.AccessToken, params.After)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error getting latest tracks"})
	}

	return c.JSON(http.StatusOK, dto.ListeningHistoryFromDomain(history))
}
