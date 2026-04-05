package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

type DonationRepository struct {
	donations map[string]*model.Donation
	mu        sync.RWMutex
	nextID    int
}

func NewDonationRepository() *DonationRepository {
	return &DonationRepository{
		donations: make(map[string]*model.Donation),
	}
}

func (r *DonationRepository) List() []*model.Donation {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.Donation, 0, len(r.donations))
	for _, d := range r.donations {
		result = append(result, d)
	}
	return result
}

func (r *DonationRepository) Create(d *model.Donation) (*model.Donation, error) {
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
