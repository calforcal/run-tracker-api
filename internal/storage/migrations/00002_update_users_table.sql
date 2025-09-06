-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
  ADD COLUMN spotify_id VARCHAR(255),
  ADD COLUMN spotify_access_token VARCHAR(500),
  ADD COLUMN spotify_refresh_token VARCHAR(500),
  ADD COLUMN spotify_expires_at BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
  DROP COLUMN spotify_expires_at,
  DROP COLUMN spotify_refresh_token,
  DROP COLUMN spotify_access_token,
  DROP COLUMN spotify_id;
-- +goose StatementEnd