package user

import (
	"fmt"
	"net/http"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/spotify"
	"run-tracker-api/internal/users"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	UserHandler struct {
		config         *config.Config
		logger         *zap.Logger
		spotifyService *spotify.SpotifyService
		userService    *users.UserService
	}

	ListeningHistoryRequest struct {
		After  int64 `query:"after"`
		Before int64 `query:"before"`
	}
)

func New(cfg *config.Config, spotifyService *spotify.SpotifyService, userService *users.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		config:         cfg,
		spotifyService: spotifyService,
		userService:    userService,
		logger:         logger,
	}
}

func (h *UserHandler) GetListeningHistory(c echo.Context) error {
	uuid := c.Get("uuid").(string)

	var params ListeningHistoryRequest
	err := c.Bind(&params)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request parameters"})
	}

	user, err := h.userService.GetUserByUUID(uuid)
	if err != nil {
		h.logger.Info(fmt.Sprintf("no user found for uuid: %s", uuid))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "error getting user"})
	}

	// if user.SpotifyExpiresAt != nil && *user.SpotifyExpiresAt < time.Now().Unix() {
	// 	h.logger.Info("token is expired")

	// 	tokenResponse, err := h.spotifyService.RefreshToken(*user.SpotifyRefreshToken)
	// 	fmt.Println("token is fucked   ", err)
	// 	if err != nil {
	// 		h.logger.Error("failed to refresh spotify token", zap.Error(err))
	// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to refresh token"})
	// 	}

	// 	fmt.Println("token is fucked  response  ", tokenResponse)

	// 	updatedUser, err := h.userService.UpdateSpotifyUser(&user, &tokenResponse)
	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error updating user"})
	// 	}
	// 	user = *updatedUser
	// }

	latestTracks, err := h.spotifyService.GetListeningHistory(*user.SpotifyAccessToken, params.After)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error getting latest tracks"})
	}

	return c.JSON(http.StatusOK, latestTracks)
}
