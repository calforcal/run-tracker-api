package spotify

import (
	"time"

	"run-tracker-api/internal/domain"
	spotifypkg "run-tracker-api/pkg/spotify"
)

func toDomainSong(t spotifypkg.TrackInfo) domain.Song {
	var artist string
	if len(t.Artists) > 0 {
		artist = t.Artists[0].Name
	}

	var imageURL string
	if len(t.Album.Images) > 0 {
		imageURL = t.Album.Images[0].URL
	}

	return domain.Song{
		Title:      t.Name,
		Artist:     artist,
		AlbumTitle: t.Album.Name,
		DurationMs: t.DurationMs,
		ImageURL:   imageURL,
		SongURI:    t.URI,
		SpotifyID:  t.ID,
	}
}

func toDomainListeningHistory(history spotifypkg.ListeningHistory) []domain.ListeningHistoryItem {
	out := make([]domain.ListeningHistoryItem, len(history.Items))
	for i, item := range history.Items {
		playedAt, _ := time.Parse(time.RFC3339, item.PlayedAt)
		out[i] = domain.ListeningHistoryItem{
			Song:     toDomainSong(item.Track),
			PlayedAt: playedAt,
		}
	}
	return out
}

// toDomainCredentials converts a Spotify OAuth token response into domain
// credentials, computing an absolute expiry from the relative ExpiresIn
// (seconds). The Storage layer previously computed this itself in two of
// three places without multiplying by time.Second, effectively persisting
// an already-expired token; computing it once, correctly, here removes
// that bug.
func toDomainCredentials(t spotifypkg.TokenResponse) domain.SpotifyCredentials {
	return domain.SpotifyCredentials{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(t.ExpiresIn) * time.Second),
	}
}
