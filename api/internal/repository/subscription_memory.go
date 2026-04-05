package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

type SubscriptionRepository struct {
	subscriptions map[string]*model.PushSubscription
	mu            sync.RWMutex
	nextID        int
}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{
		subscriptions: make(map[string]*model.PushSubscription),
	}
}

func (r *SubscriptionRepository) List() []*model.PushSubscription {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.PushSubscription, 0, len(r.subscriptions))
	for _, s := range r.subscriptions {
		result = append(result, s)
	}
	return result
}

func (r *SubscriptionRepository) Create(s *model.PushSubscription) (*model.PushSubscription, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	s.ID = fmt.Sprintf("sub-%d", r.nextID)
	s.CreatedAt = time.Now()

	r.subscriptions[s.ID] = s
	return s, nil
}

func (r *SubscriptionRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.subscriptions[id]; !ok {
		return fmt.Errorf("subscription not found: %s", id)
	}
	delete(r.subscriptions, id)
	return nil
}

func (r *SubscriptionRepository) DeleteByEndpoint(endpoint string) error {
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
