package auth

import (
	"net/http"
	"strings"

	"run-tracker-api/api/dto"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/ports"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthHandler struct {
	config      *config.Config
	userService ports.UserService
	authService ports.AuthService
	logger      *zap.Logger
}

func New(cfg *config.Config, userService ports.UserService, authService ports.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		config:      cfg,
		userService: userService,
		authService: authService,
		logger:      logger,
	}
}

// Login exchanges a Strava OAuth code, upserts the user, ensures a linked
// Spotify token is still valid, and issues a session token.
func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.ExchangeCodeRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Info("missing code from request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := h.userService.LoginWithStrava(ctx, req.Code)
	if err != nil {
		h.logger.Info("failed to log in with strava", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to exchange code for token"})
	}

	token, err := h.authService.IssueToken(ctx, user)
	if err != nil {
		h.logger.Info("failed to issue new token", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to authorize user"})
	}

	return c.JSON(http.StatusOK, dto.TokenResponse{AccessToken: token})
}

// AuthorizeStravaUser exchanges a Strava OAuth code, upserts the user, and
// issues a session token, without touching any linked Spotify credentials.
func (h *AuthHandler) AuthorizeStravaUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.ExchangeCodeRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Info("missing code from request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := h.userService.ExchangeStravaCode(ctx, req.Code)
	if err != nil {
		h.logger.Info("failed to exchange code for token", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to exchange code for token"})
	}

	token, err := h.authService.IssueToken(ctx, user)
	if err != nil {
		h.logger.Info("failed to issue new token", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to authorize user"})
	}

	return c.JSON(http.StatusOK, dto.TokenResponse{AccessToken: token})
}

// AuthorizeSpotifyUser links a Spotify account to the already-authenticated
// (Strava) user identified by the bearer token on the request.
func (h *AuthHandler) AuthorizeSpotifyUser(c echo.Context) error {
	ctx := c.Request().Context()

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "unauthorized spotify login attempt"})
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.authService.ParseToken(ctx, tokenStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
	}
	if claims.UUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "No valid user UUID provided"})
	}

	var req dto.ExchangeCodeRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Info("missing code from request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := h.userService.LinkSpotifyAccount(ctx, claims.UUID, req.Code, h.config.SpotifyRedirectURI)
	if err != nil {
		h.logger.Info("failed to link spotify account", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update user"})
	}

	token, err := h.authService.IssueToken(ctx, user)
	if err != nil {
		h.logger.Info("failed to issue new token", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to authorize user"})
	}

	return c.JSON(http.StatusOK, dto.TokenResponse{AccessToken: token})
}
