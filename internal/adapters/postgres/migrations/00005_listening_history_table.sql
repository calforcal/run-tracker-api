-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    album_title VARCHAR(255) NOT NULL,
    duration INTEGER NOT NULL,
    image_url VARCHAR(255),
    song_uri VARCHAR(255),
    spotify_id VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_activity_songs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    activity_id BIGINT NOT NULL,
    song_id BIGINT REFERENCES songs(id) ON DELETE CASCADE,
    played_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_activity ON user_activity_songs(user_id, activity_id, song_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_activity_songs;
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd