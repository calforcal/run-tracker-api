package auth

import "github.com/golang-jwt/jwt/v5"

type (
	CustomClaims struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
		jwt.RegisteredClaims
	}
)
