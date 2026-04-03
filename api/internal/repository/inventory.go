package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

type InventoryRepository struct {
	mu        sync.RWMutex
	inventory map[model.BloodType]*model.BloodInventory
	nextID    int
}

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

func (r *InventoryRepository) List() []*model.BloodInventory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.BloodInventory, 0, len(r.inventory))
	for _, inv := range r.inventory {
		result = append(result, inv)
	}
	return result
}

func (r *InventoryRepository) Update(bloodType model.BloodType, level model.InventoryLevel) (*model.BloodInventory, error) {
	if !model.IsValidBloodType(bloodType) {
		return nil, fmt.Errorf("invalid blood type: %s", bloodType)
	}
	if !model.IsValidInventoryLevel(level) {
		return nil, fmt.Errorf("invalid inventory level: %s", level)
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

func (r *InventoryRepository) GetByBloodType(bt model.BloodType) (*model.BloodInventory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	inv, ok := r.inventory[bt]
	if !ok {
		return nil, fmt.Errorf("inventory not found for blood type: %s", bt)
	}
	return inv, nil
}
