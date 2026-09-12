-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS webhook_subscriptions (
  id SERIAL PRIMARY KEY,
  strava_id BIGINT NOT NULL,
  callback_url VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS webhook_subscriptions;
-- +goose StatementEnd