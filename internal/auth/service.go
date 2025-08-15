package auth

import (
	"fmt"
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
	fmt.Printf("Parsing token: %s\n", tokenStr)
	fmt.Printf("Using secret: %v\n", s.Config.JwtSecret)

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Convert string to []byte
		return []byte(s.Config.JwtSecret), nil
	})

	fmt.Printf("Parse error: %v\n", err)
	fmt.Printf("Token valid: %v\n", token.Valid)

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token validation failed: %v", err)
}
