package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"run-tracker-api/internal/domain"
	"run-tracker-api/internal/ports"
)

type webhookRepository struct {
	db *sql.DB
}

// NewWebhookRepository builds a ports.WebhookRepository backed by Postgres.
func NewWebhookRepository(db *sql.DB) ports.WebhookRepository {
	return &webhookRepository{db: db}
}

var _ ports.WebhookRepository = (*webhookRepository)(nil)

func (r *webhookRepository) Create(ctx context.Context, sub domain.WebhookSubscription) (domain.WebhookSubscription, error) {
	query := `INSERT INTO webhook_subscriptions (strava_id, callback_url) VALUES ($1, $2) RETURNING id, strava_id, callback_url`

	var result domain.WebhookSubscription
	err := r.db.QueryRowContext(ctx, query, sub.StravaID, sub.CallbackURL).Scan(&result.ID, &result.StravaID, &result.CallbackURL)
	if err != nil {
		return domain.WebhookSubscription{}, fmt.Errorf("error writing webhook subscription to database: %w", err)
	}
	return result, nil
}

func (r *webhookRepository) Get(ctx context.Context) (domain.WebhookSubscription, error) {
	query := `SELECT id, strava_id, callback_url FROM webhook_subscriptions`

	var result domain.WebhookSubscription
	err := r.db.QueryRowContext(ctx, query).Scan(&result.ID, &result.StravaID, &result.CallbackURL)
	if err != nil {
		return domain.WebhookSubscription{}, fmt.Errorf("error reading webhook subscription from database: %w", err)
	}
	return result, nil
}

func (r *webhookRepository) Delete(ctx context.Context, stravaSubscriptionID int) error {
	query := `DELETE FROM webhook_subscriptions WHERE strava_id = $1`
	if _, err := r.db.ExecContext(ctx, query, stravaSubscriptionID); err != nil {
		return fmt.Errorf("error deleting webhook subscription: %w", err)
	}
	return nil
}
