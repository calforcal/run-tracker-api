// Package auth implements ports.AuthService: issuing and parsing this
// application's own session JWTs. This is a pure crypto/algorithm concern
// (not an external API or a database), so unlike the strava/spotify/user
// services it needs no adapter of its own.
package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"

	"github.com/golang-jwt/jwt/v5"
)

// jwtClaims is the wire shape signed into the token. It's kept private and
// separate from domain.Claims so the JWT's exact on-the-wire shape (a
// space-joined "scopes" string, for existing frontend compatibility) is an
// implementation detail of this package, not a domain concept.
type jwtClaims struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Scopes string `json:"scopes"`
	jwt.RegisteredClaims
}

type service struct {
	jwtSecret string
}

// New builds a ports.AuthService.
func New(jwtSecret string) ports.AuthService {
	return &service{jwtSecret: jwtSecret}
}

var _ ports.AuthService = (*service)(nil)

func (s *service) IssueToken(ctx context.Context, user domain.User) (string, error) {
	var scopes string
	switch {
	case user.Strava.AthleteID != 0 && user.Spotify != nil:
		scopes = "strava spotify"
	case user.Strava.AthleteID != 0:
		scopes = "strava"
	default:
		return "", fmt.Errorf("invalid user state")
	}

	claims := jwtClaims{
		UUID:   user.UUID,
		Name:   user.Name,
		Scopes: scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "run-tracker",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *service) ParseToken(ctx context.Context, tokenStr string) (domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return domain.Claims{}, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return domain.Claims{}, fmt.Errorf("token validation failed")
	}

	result := domain.Claims{
		UUID:   claims.UUID,
		Name:   claims.Name,
		Scopes: strings.Fields(claims.Scopes),
	}
	if claims.IssuedAt != nil {
		result.IssuedAt = claims.IssuedAt.Time
	}
	if claims.ExpiresAt != nil {
		result.ExpiresAt = claims.ExpiresAt.Time
	}

	return result, nil
}
