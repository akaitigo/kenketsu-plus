package repository_test

import (
	"context"
	"testing"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
)

// TestSubscriptionRepository_UpsertByEndpoint verifies the UPSERT semantics (#18):
// re-subscribing with the same endpoint updates the existing record in place
// instead of creating a duplicate that would cause duplicate notifications.
func TestSubscriptionRepository_UpsertByEndpoint(t *testing.T) {
	t.Parallel()

	repo := repository.NewSubscriptionRepository()
	ctx := context.Background()

	first, err := repo.Create(ctx, &model.PushSubscription{
		Endpoint: "https://push.example.com/sub-a",
		P256dh:   "key-1",
		Auth:     "auth-1",
	})
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	second, err := repo.Create(ctx, &model.PushSubscription{
		Endpoint: "https://push.example.com/sub-a",
		P256dh:   "key-2",
		Auth:     "auth-2",
	})
	if err != nil {
		t.Fatalf("second create failed: %v", err)
	}

	subs := repo.List(ctx)
	if len(subs) != 1 {
		t.Fatalf("expected 1 subscription after re-subscribe, got %d", len(subs))
	}
	if second.ID != first.ID {
		t.Errorf("expected stable ID %q across upsert, got %q", first.ID, second.ID)
	}
	if subs[0].P256dh != "key-2" || subs[0].Auth != "auth-2" {
		t.Errorf("expected keys updated to key-2/auth-2, got %s/%s", subs[0].P256dh, subs[0].Auth)
	}
}

// TestSubscriptionRepository_DistinctEndpoints verifies distinct endpoints are
// still stored as separate records.
func TestSubscriptionRepository_DistinctEndpoints(t *testing.T) {
	t.Parallel()

	repo := repository.NewSubscriptionRepository()
	ctx := context.Background()

	if _, err := repo.Create(ctx, &model.PushSubscription{Endpoint: "https://push.example.com/a", P256dh: "k", Auth: "a"}); err != nil {
		t.Fatalf("create a failed: %v", err)
	}
	if _, err := repo.Create(ctx, &model.PushSubscription{Endpoint: "https://push.example.com/b", P256dh: "k", Auth: "a"}); err != nil {
		t.Fatalf("create b failed: %v", err)
	}

	if subs := repo.List(ctx); len(subs) != 2 {
		t.Fatalf("expected 2 subscriptions for distinct endpoints, got %d", len(subs))
	}
}
