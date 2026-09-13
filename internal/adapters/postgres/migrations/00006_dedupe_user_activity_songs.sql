-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_activity_songs
  ADD CONSTRAINT unique_user_song_played_at UNIQUE (user_id, song_id, played_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_activity_songs
  DROP CONSTRAINT unique_user_song_played_at;
-- +goose StatementEnd
