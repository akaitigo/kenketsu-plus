package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

// InventoryRepository is the in-memory implementation of InventoryRepo.
type InventoryRepository struct {
	inventory map[model.BloodType]*model.BloodInventory
	mu        sync.RWMutex
	nextID    int
}

// NewInventoryRepository creates a new in-memory inventory repository with default entries.
func NewInventoryRepository() *InventoryRepository {
	repo := &InventoryRepository{
		inventory: make(map[model.BloodType]*model.BloodInventory),
	}
	for _, bt := range model.ValidBloodTypes {
		repo.nextID++
		repo.inventory[bt] = &model.BloodInventory{
			ID:        fmt.Sprintf("inv-%d", repo.nextID),
			BloodType: bt,
			Level:     model.InventoryLevelNormal,
			UpdatedAt: time.Now(),
		}
	}
	return repo
}

// List returns all blood inventory records from memory.
func (r *InventoryRepository) List(_ context.Context) []*model.BloodInventory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.BloodInventory, 0, len(r.inventory))
	for _, inv := range r.inventory {
		result = append(result, inv)
	}
	return result
}

// Update changes the inventory level for the specified blood type.
func (r *InventoryRepository) Update(_ context.Context, bloodType model.BloodType, level model.InventoryLevel) (*model.BloodInventory, error) {
	if !model.IsValidBloodType(bloodType) {
		return nil, model.ErrFieldInvalid("bloodType", "invalid blood type: "+string(bloodType))
	}
	if !model.IsValidInventoryLevel(level) {
		return nil, model.ErrFieldInvalid("level", "invalid inventory level: "+string(level))
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	inv, ok := r.inventory[bloodType]
	if !ok {
		return nil, fmt.Errorf("inventory not found for blood type: %s", bloodType)
	}

	inv.Level = level
	inv.UpdatedAt = time.Now()
	return inv, nil
}

// GetByBloodType returns the inventory record for the specified blood type.
func (r *InventoryRepository) GetByBloodType(_ context.Context, bt model.BloodType) (*model.BloodInventory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	inv, ok := r.inventory[bt]
	if !ok {
		return nil, fmt.Errorf("inventory not found for blood type: %s", bt)
	}
	return inv, nil
}

// Compile-time interface check.
var _ InventoryRepo = (*InventoryRepository)(nil)
