package domain

import "time"

// Song is a track that was played, persisted independently of any single
// listen so repeated listens reuse the same row.
type Song struct {
	ID         int
	Title      string
	Artist     string
	AlbumTitle string
	DurationMs int
	ImageURL   string
	SongURI    string
	SpotifyID  string
}

// ListeningHistoryItem is a single play of a Song.
type ListeningHistoryItem struct {
	Song     Song
	PlayedAt time.Time
}
