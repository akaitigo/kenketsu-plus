package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

// SubscriptionRepository is the in-memory implementation of SubscriptionRepo.
type SubscriptionRepository struct {
	subscriptions map[string]*model.PushSubscription
	mu            sync.RWMutex
	nextID        int
}

// NewSubscriptionRepository creates a new in-memory subscription repository.
func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{
		subscriptions: make(map[string]*model.PushSubscription),
	}
}

// List returns all push subscriptions from memory.
func (r *SubscriptionRepository) List(_ context.Context) []*model.PushSubscription {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.PushSubscription, 0, len(r.subscriptions))
	for _, s := range r.subscriptions {
		result = append(result, s)
	}
	return result
}

// Create inserts a new push subscription into memory, or updates the keys of an
// existing one when the endpoint is already registered (UPSERT). This mirrors
// the PostgreSQL UNIQUE(endpoint) behavior so that re-subscribing does not create
// duplicate records that would cause duplicate notifications (#18).
func (r *SubscriptionRepository) Create(_ context.Context, s *model.PushSubscription) (*model.PushSubscription, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.subscriptions {
		if existing.Endpoint == s.Endpoint {
			existing.P256dh = s.P256dh
			existing.Auth = s.Auth
			s.ID = existing.ID
			s.CreatedAt = existing.CreatedAt
			return existing, nil
		}
	}

	r.nextID++
	s.ID = fmt.Sprintf("sub-%d", r.nextID)
	s.CreatedAt = time.Now()

	r.subscriptions[s.ID] = s
	return s, nil
}

// Delete removes a push subscription by its ID.
func (r *SubscriptionRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.subscriptions[id]; !ok {
		return fmt.Errorf("subscription not found: %s", id)
	}
	delete(r.subscriptions, id)
	return nil
}

// DeleteByEndpoint removes a push subscription by its endpoint URL.
func (r *SubscriptionRepository) DeleteByEndpoint(_ context.Context, endpoint string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, s := range r.subscriptions {
		if s.Endpoint == endpoint {
			delete(r.subscriptions, id)
			return nil
		}
	}
	return fmt.Errorf("subscription not found for endpoint")
}

// Compile-time interface check.
var _ SubscriptionRepo = (*SubscriptionRepository)(nil)
