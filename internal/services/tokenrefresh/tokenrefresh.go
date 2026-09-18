// Package tokenrefresh centralizes the "refresh this OAuth token if it's
// close to expiring" logic that used to be duplicated (and always
// unconditional) across the Login handler, the webhook processor, and a
// dead attempt in the listening-history handler. It's a shared helper
// package, not a service - services/user and services/webhook both call
// into it, but neither depends on the other.
package tokenrefresh

import (
	"context"
	"time"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"
)

// skew is how far ahead of actual expiry we proactively refresh.
const skew = 5 * time.Minute

// EnsureValidStravaToken returns user unchanged if its Strava access token
// still has more than skew left, otherwise refreshes it and persists the
// result.
func EnsureValidStravaToken(ctx context.Context, user domain.User, repo ports.UserRepository, provider ports.StravaProvider) (domain.User, error) {
	if time.Now().Add(skew).Before(user.Strava.ExpiresAt) {
		return user, nil
	}

	creds, err := provider.RefreshToken(ctx, user.Strava.RefreshToken)
	if err != nil {
		return domain.User{}, err
	}

	return repo.UpdateStravaCredentials(ctx, user.Strava.AthleteID, creds)
}

// EnsureValidSpotifyToken returns user unchanged if it has no linked
// Spotify account, or if its Spotify access token still has more than skew
// left. Otherwise it refreshes the token and persists the result.
func EnsureValidSpotifyToken(ctx context.Context, user domain.User, repo ports.UserRepository, provider ports.SpotifyProvider) (domain.User, error) {
	if user.Spotify == nil {
		return user, nil
	}
	if time.Now().Add(skew).Before(user.Spotify.ExpiresAt) {
		return user, nil
	}

	creds, err := provider.RefreshToken(ctx, user.Spotify.RefreshToken)
	if err != nil {
		return domain.User{}, err
	}
	// The refresh response doesn't carry the Spotify user id back.
	creds.SpotifyID = user.Spotify.SpotifyID
	// Spotify doesn't always return a new refresh_token on refresh - when a
	// client's token hasn't been rotated, the field is simply omitted. The
	// existing refresh token is still valid in that case and must be kept;
	// overwriting it with the resulting empty string would permanently break
	// every future refresh for this user.
	if creds.RefreshToken == "" {
		creds.RefreshToken = user.Spotify.RefreshToken
	}

	return repo.UpdateSpotifyCredentials(ctx, user.Spotify.SpotifyID, creds)
}
