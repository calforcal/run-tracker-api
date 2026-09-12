package domain

import "time"

// Claims is the set of identity/authorization data carried by an issued
// session token.
type Claims struct {
	UUID      string
	Name      string
	Scopes    []string
	IssuedAt  time.Time
	ExpiresAt time.Time
}
