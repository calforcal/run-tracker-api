package auth

import (
	"errors"
	"net/http"
	"run-tracker-api/internal/auth"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/strava"
	"run-tracker-api/internal/users"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	AuthHandler struct {
		config        *config.Config
		logger        *zap.Logger
		stravaService *strava.StravaService
		userService   *users.UserService
		authService   *auth.AuthService
	}

	ExchangeCodeForTokenRequest struct {
		Code string `json:"code"`
	}

	Token struct {
		AccessToken string `json:"access_token"`
	}
)

func New(cfg *config.Config, stravaService *strava.StravaService, userService *users.UserService, authService *auth.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		config:        cfg,
		stravaService: stravaService,
		userService:   userService,
		authService:   authService,
		logger:        logger,
	}
}

func (h *AuthHandler) AuthorizeStravaUser(c echo.Context) error {
	var req ExchangeCodeForTokenRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Info("missing code from request: %d", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	tokenResponse, err := h.exchangeCodeForToken(req)
	if err != nil {
		h.logger.Info("failed to exchange code for token: %d", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to exchange code for token"})
	}

	user, err := h.userService.CreateOrUpdateUser(tokenResponse)
	if err != nil {
		h.logger.Info("failed to create or update user: %d", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "error authorizing user")
	}

	expiresAt := time.Unix(int64(user.StravaExpiresAt), 0)
	if time.Now().UTC().After(expiresAt) {
		refreshResponse, err := h.stravaService.RefreshToken(user.StravaRefreshToken)
		if err != nil {
			h.logger.Info("failed to refresh users token: %d", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, "error authorizing user")
		}
		stravaToken := strava.TokenResponse{
			TokenType:    "Bearer",
			AccessToken:  refreshResponse.AccessToken,
			RefreshToken: refreshResponse.RefreshToken,
			ExpiresAt:    refreshResponse.ExpiresAt,
			ExpiresIn:    refreshResponse.ExpiresIn,
			Athlete:      tokenResponse.Athlete,
		}
		user, err := h.userService.CreateOrUpdateUser(&stravaToken)
		if err != nil {
			h.logger.Info("failed to update users token in database: %d", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, "error authorizing user")
		}
		token, err := h.authService.IssueJwt(user)
		if err != nil {
			h.logger.Info("failed to issue new jwt: %d", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to authorize user"})
		}

		return c.JSON(http.StatusOK, Token{AccessToken: token})
	}

	token, err := h.authService.IssueJwt(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to authorize user"})
	}

	return c.JSON(http.StatusOK, Token{AccessToken: token})
}

func (h *AuthHandler) exchangeCodeForToken(req ExchangeCodeForTokenRequest) (*strava.TokenResponse, error) {
	tokenResponse, err := h.stravaService.ExchangeCodeForToken(req.Code)

	if err != nil {
		return &strava.TokenResponse{}, err
	}

	if tokenResponse.Athlete.ID <= 0 {
		return &strava.TokenResponse{}, errors.New("error authenticating user")
	}

	return &tokenResponse, nil
}
