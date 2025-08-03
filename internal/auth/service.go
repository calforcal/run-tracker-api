package auth

import (
	"run-tracker-api/internal/config"
	"run-tracker-api/internal/storage"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type (
	AuthService struct {
		Config *config.Config
		Logger *zap.Logger
	}
)

func New(cfg *config.Config, logger *zap.Logger) *AuthService {
	return &AuthService{
		Config: cfg,
		Logger: logger,
	}
}

func (s *AuthService) IssueJwt(user *storage.User) (string, error) {
	secret := s.Config.JwtSecret

	claims := CustomClaims{
		UUID: user.UUID,
		Name: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "run-tracker",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *AuthService) ParseJWT(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.Config.JwtSecret, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
