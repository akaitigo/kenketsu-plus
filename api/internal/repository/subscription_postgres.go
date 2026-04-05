package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

// PgSubscriptionRepository is the PostgreSQL implementation of SubscriptionRepo.
type PgSubscriptionRepository struct {
	db *sql.DB
}

// NewPgSubscriptionRepository creates a new PostgreSQL-backed subscription repository.
func NewPgSubscriptionRepository(db *sql.DB) *PgSubscriptionRepository {
	return &PgSubscriptionRepository{db: db}
}

// List returns all push subscriptions.
func (r *PgSubscriptionRepository) List(ctx context.Context) []*model.PushSubscription {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, endpoint, p256dh, auth, created_at
		FROM push_subscriptions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil
	}
	defer func() { _ = rows.Close() }()

	var subs []*model.PushSubscription
	for rows.Next() {
		var s model.PushSubscription
		if err := rows.Scan(&s.ID, &s.Endpoint, &s.P256dh, &s.Auth, &s.CreatedAt); err != nil {
			continue
		}
		subs = append(subs, &s)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	if subs == nil {
		return []*model.PushSubscription{}
	}
	return subs
}

// Create inserts a new push subscription.
func (r *PgSubscriptionRepository) Create(ctx context.Context, s *model.PushSubscription) (*model.PushSubscription, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO push_subscriptions (endpoint, p256dh, auth)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, s.Endpoint, s.P256dh, s.Auth).Scan(&s.ID, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Delete removes a push subscription by ID.
func (r *PgSubscriptionRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM push_subscriptions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("subscription not found: %s", id)
	}
	return nil
}

// DeleteByEndpoint removes a push subscription by endpoint URL.
func (r *PgSubscriptionRepository) DeleteByEndpoint(ctx context.Context, endpoint string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM push_subscriptions WHERE endpoint = $1`, endpoint)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("subscription not found for endpoint")
	}
	return nil
}

// Compile-time interface check.
var _ SubscriptionRepo = (*PgSubscriptionRepository)(nil)
