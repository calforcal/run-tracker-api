package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"
)

type listeningHistoryRepository struct {
	db *sql.DB
}

// NewListeningHistoryRepository builds a ports.ListeningHistoryRepository
// backed by Postgres.
func NewListeningHistoryRepository(db *sql.DB) ports.ListeningHistoryRepository {
	return &listeningHistoryRepository{db: db}
}

var _ ports.ListeningHistoryRepository = (*listeningHistoryRepository)(nil)

func (r *listeningHistoryRepository) SaveEntry(ctx context.Context, userID int, activityID int, item domain.ListeningHistoryItem) error {
	song, err := r.getOrCreateSong(ctx, item.Song)
	if err != nil {
		return fmt.Errorf("error creating song in database: %w", err)
	}

	if err := r.saveUserSong(ctx, userID, activityID, song.ID, item.PlayedAt); err != nil {
		return fmt.Errorf("error creating song:user association: %w", err)
	}

	return nil
}

func (r *listeningHistoryRepository) GetForActivity(ctx context.Context, userID, activityID int) ([]domain.ListeningHistoryItem, error) {
	query := `
		SELECT s.title, s.artist, s.album_title, s.duration, s.image_url, s.song_uri, s.spotify_id, uas.played_at
		FROM user_activity_songs uas
		JOIN songs s ON s.id = uas.song_id
		WHERE uas.user_id = $1 AND uas.activity_id = $2
		ORDER BY uas.played_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, activityID)
	if err != nil {
		return nil, fmt.Errorf("error reading listening history for activity: %w", err)
	}
	defer rows.Close()

	items := make([]domain.ListeningHistoryItem, 0)
	for rows.Next() {
		var item domain.ListeningHistoryItem
		if err := rows.Scan(
			&item.Song.Title, &item.Song.Artist, &item.Song.AlbumTitle, &item.Song.DurationMs,
			&item.Song.ImageURL, &item.Song.SongURI, &item.Song.SpotifyID, &item.PlayedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning listening history row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating listening history rows: %w", err)
	}

	return items, nil
}

func (r *listeningHistoryRepository) getOrCreateSong(ctx context.Context, song domain.Song) (domain.Song, error) {
	query := `
		INSERT INTO songs (title, artist, album_title, duration, image_url, song_uri, spotify_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (spotify_id) DO UPDATE SET
			title = EXCLUDED.title,
			artist = EXCLUDED.artist,
			album_title = EXCLUDED.album_title,
			duration = EXCLUDED.duration,
			image_url = EXCLUDED.image_url,
			song_uri = EXCLUDED.song_uri
		RETURNING id, title, artist, album_title, duration, image_url, song_uri, spotify_id
	`

	var result domain.Song
	err := r.db.QueryRowContext(ctx, query,
		song.Title, song.Artist, song.AlbumTitle, song.DurationMs, song.ImageURL, song.SongURI, song.SpotifyID,
	).Scan(
		&result.ID, &result.Title, &result.Artist, &result.AlbumTitle, &result.DurationMs, &result.ImageURL, &result.SongURI, &result.SpotifyID,
	)
	if err != nil {
		return domain.Song{}, fmt.Errorf("failed to get or create song: %w", err)
	}

	return result, nil
}

func (r *listeningHistoryRepository) saveUserSong(ctx context.Context, userID, activityID, songID int, playedAt time.Time) error {
	query := `
		INSERT INTO user_activity_songs (user_id, activity_id, song_id, played_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, song_id, played_at) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, activityID, songID, playedAt)
	if err != nil {
		return fmt.Errorf("error writing user song to database: %w", err)
	}
	return nil
}
