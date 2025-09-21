-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN strava_id SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT unique_strava_id UNIQUE (strava_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT unique_strava_id;
ALTER TABLE users ALTER COLUMN strava_id DROP NOT NULL;
-- +goose StatementEnd