package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

// DonationRepository is the in-memory implementation of DonationRepo.
type DonationRepository struct {
	donations map[string]*model.Donation
	mu        sync.RWMutex
	nextID    int
}

// NewDonationRepository creates a new in-memory donation repository.
func NewDonationRepository() *DonationRepository {
	return &DonationRepository{
		donations: make(map[string]*model.Donation),
	}
}

// List returns all donation records from memory.
func (r *DonationRepository) List(_ context.Context) []*model.Donation {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.Donation, 0, len(r.donations))
	for _, d := range r.donations {
		result = append(result, d)
	}
	return result
}

// Create inserts a new donation record into memory.
func (r *DonationRepository) Create(_ context.Context, d *model.Donation) (*model.Donation, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	d.ID = fmt.Sprintf("donation-%d", r.nextID)
	d.CreatedAt = time.Now()

	r.donations[d.ID] = d
	return d, nil
}

// Compile-time interface check.
var _ DonationRepo = (*DonationRepository)(nil)
