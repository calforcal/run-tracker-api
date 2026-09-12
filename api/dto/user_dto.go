package dto

import "run-tracker-api/internal/domain"

// ListeningHistoryRequest binds query params for GET /api/users/listening-history.
type ListeningHistoryRequest struct {
	After  int64 `query:"after"`
	Before int64 `query:"before"`
}

// SongResponse is the persisted (already-flattened, single artist/image)
// shape of a listened-to track. See domain.Song for why this carries one
// artist/image rather than the full arrays Spotify's API returns - it's
// the same shape saved to the database, reused here for the API response.
type SongResponse struct {
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	AlbumTitle string `json:"album_title"`
	DurationMs int    `json:"duration_ms"`
	ImageURL   string `json:"image_url"`
	URI        string `json:"uri"`
	SpotifyID  string `json:"spotify_id"`
}

type ListeningHistoryItemResponse struct {
	Song     SongResponse `json:"song"`
	PlayedAt string       `json:"played_at"`
}

type ListeningHistoryResponse struct {
	Items []ListeningHistoryItemResponse `json:"items"`
}

func ListeningHistoryFromDomain(items []domain.ListeningHistoryItem) ListeningHistoryResponse {
	out := make([]ListeningHistoryItemResponse, len(items))
	for i, item := range items {
		out[i] = ListeningHistoryItemResponse{
			Song: SongResponse{
				Title:      item.Song.Title,
				Artist:     item.Song.Artist,
				AlbumTitle: item.Song.AlbumTitle,
				DurationMs: item.Song.DurationMs,
				ImageURL:   item.Song.ImageURL,
				URI:        item.Song.SongURI,
				SpotifyID:  item.Song.SpotifyID,
			},
			PlayedAt: formatTime(item.PlayedAt),
		}
	}
	return ListeningHistoryResponse{Items: out}
}
