package auth

import (
	"errors"
	"fmt"
	"net/http"
	"run-tracker-api/internal/auth"
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/spotify"
	"run-tracker-api/internal/storage"
	"run-tracker-api/internal/strava"
	"run-tracker-api/internal/users"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	AuthHandler struct {
		config         *config.Config
		logger         *zap.Logger
		stravaService  *strava.StravaService
		spotifyService *spotify.SpotifyService
		userService    *users.UserService
		authService    *auth.AuthService
	}

	ExchangeCodeForTokenRequest struct {
		Code string `json:"code"`
	}

	Token struct {
		AccessToken string `json:"access_token"`
	}
)

func New(cfg *config.Config, stravaService *strava.StravaService, spotifyService *spotify.SpotifyService, userService *users.UserService, authService *auth.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		config:         cfg,
		stravaService:  stravaService,
		spotifyService: spotifyService,
		userService:    userService,
		authService:    authService,
		logger:         logger,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req ExchangeCodeForTokenRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Info("missing code from request: %s", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	tokenResponse, err := h.exchangeCodeForToken(req)
	if err != nil {
		h.logger.Info("failed to exchange code for token: %d", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to exchange code for token"})
	}

	user, err := h.userService.CreateOrUpdateUser(tokenResponse)
	if err != nil {
		h.logger.Info("failed to create or upsert user: %d", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "error authorizing user")
	}

	// Refresh Spotify token
	if user.SpotifyID != nil {
		tokenResponse, err := h.spotifyService.RefreshToken(*user.SpotifyRefreshToken)
		fmt.Println("token is fucked   ", err)
		if err != nil {
			h.logger.Error("failed to refresh spotify token", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to refresh token"})
		}

		fmt.Println("token is fucked  response  ", tokenResponse)

		updatedUser, err := h.userService.UpdateSpotifyUser(user, &tokenResponse)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error updating user"})
		}
		user = updatedUser
	}

	fmt.Println("USIE", user)
	token, err := h.authService.IssueJwt(user)
	if err != nil {
		h.logger.Info("failed to issue new jwt: %d", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to authorize user"})
	}

	return c.JSON(http.StatusOK, Token{AccessToken: token})
}

func (h *AuthHandler) AuthorizeStravaUser(c echo.Context) error {
	var req ExchangeCodeForTokenRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Info("missing code from request: %s", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	tokenResponse, err := h.exchangeCodeForToken(req)
	if err != nil {
		h.logger.Info("failed to exchange code for token: %d", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to exchange code for token"})
	}

	user, err := h.userService.CreateOrUpdateUser(tokenResponse)
	if err != nil {
		h.logger.Info("failed to create or upsert user: %d", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "error authorizing user")
	}

	// Should never return here without a fresh token
	// stravaExpiresAt := int64(user.StravaExpiresAt)
	// expiresAt := time.Unix(stravaExpiresAt, 0)
	// if time.Now().UTC().After(expiresAt) {
	// refreshResponse, err := h.stravaService.RefreshToken(user.StravaRefreshToken)
	// if err != nil {
	// 	h.logger.Info("failed to refresh users token: %d", zap.Error(err))
	// 	return c.JSON(http.StatusInternalServerError, "error authorizing user")
	// }
	// stravaToken := strava.TokenResponse{
	// 	TokenType:    "Bearer",
	// 	AccessToken:  refreshResponse.AccessToken,
	// 	RefreshToken: refreshResponse.RefreshToken,
	// 	ExpiresAt:    refreshResponse.ExpiresAt,
	// 	ExpiresIn:    refreshResponse.ExpiresIn,
	// 	Athlete:      tokenResponse.Athlete,
	// }
	// user, err := h.userService.CreateOrUpdateUser(&stravaToken)
	// if err != nil {
	// 	h.logger.Info("failed to update users token in database: %d", zap.Error(err))
	// 	return c.JSON(http.StatusInternalServerError, "error authorizing user")
	// }
	token, err := h.authService.IssueJwt(user)
	if err != nil {
		h.logger.Info("failed to issue new jwt: %d", zap.Error(err))
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

func (h *AuthHandler) AuthorizeSpotifyUser(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	var uuid *string
	fmt.Println("TOKEN", token)
	if token != "" {
		tokenStr := strings.TrimPrefix(token, "Bearer ")
		claims, err := h.authService.ParseJWT(tokenStr)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
		}
		uuid = &claims.UUID
	}

	if token == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "unauthorized spotify login attempt"})
	}

	var req ExchangeCodeForTokenRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Info("missing code from request: %d", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	redirectURI := "http://127.0.0.1:5173/auth/callback/spotify"
	grantType := "authorization_code"
	tokenResponse, err := h.exchangeSpotifyCodeForToken(req, redirectURI, grantType)
	if err != nil {
		h.logger.Info("failed to exchange code for token: %v", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to exchange code for token"})
	}

	fmt.Println("TOKEN RESPONSE", tokenResponse)

	var user *storage.User
	if uuid != nil && *uuid != "" {
		user, err = h.userService.AddSpotifyToStravaUser(*uuid, tokenResponse)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update user"})
		}
	} else {
		// If no UUID provided, we can't proceed
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "No valid user UUID provided"})
	}

	fmt.Println("USER", user)

	token, err = h.authService.IssueJwt(user)
	if err != nil {
		h.logger.Info("failed to issue new jwt: %d", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to authorize user"})
	}

	return c.JSON(http.StatusOK, Token{AccessToken: token})
}

func (h *AuthHandler) exchangeSpotifyCodeForToken(req ExchangeCodeForTokenRequest, redirectURI string, grantType string) (*spotify.TokenResponse, error) {
	tokenResponse, err := h.spotifyService.ExchangeCodeForToken(req.Code, redirectURI, grantType)

	if err != nil {
		return &spotify.TokenResponse{}, err
	}

	return &tokenResponse, nil
}
