package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

// PgInventoryRepository is the PostgreSQL implementation of InventoryRepo.
type PgInventoryRepository struct {
	db *sql.DB
}

// NewPgInventoryRepository creates a new PostgreSQL-backed inventory repository.
func NewPgInventoryRepository(db *sql.DB) *PgInventoryRepository {
	return &PgInventoryRepository{db: db}
}

// List returns all blood inventory records.
func (r *PgInventoryRepository) List(ctx context.Context) []*model.BloodInventory {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, blood_type, level, updated_at
		FROM blood_inventory
		ORDER BY blood_type
	`)
	if err != nil {
		return nil
	}
	defer func() { _ = rows.Close() }()

	var inventories []*model.BloodInventory
	for rows.Next() {
		var inv model.BloodInventory
		if err := rows.Scan(&inv.ID, &inv.BloodType, &inv.Level, &inv.UpdatedAt); err != nil {
			continue
		}
		inventories = append(inventories, &inv)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	if inventories == nil {
		return []*model.BloodInventory{}
	}
	return inventories
}

// Update updates the inventory level for the specified blood type.
func (r *PgInventoryRepository) Update(ctx context.Context, bloodType model.BloodType, level model.InventoryLevel) (*model.BloodInventory, error) {
	if !model.IsValidBloodType(bloodType) {
		return nil, model.ErrFieldInvalid("bloodType", "invalid blood type: "+string(bloodType))
	}
	if !model.IsValidInventoryLevel(level) {
		return nil, model.ErrFieldInvalid("level", "invalid inventory level: "+string(level))
	}

	var inv model.BloodInventory
	err := r.db.QueryRowContext(ctx, `
		UPDATE blood_inventory SET level = $1, updated_at = NOW()
		WHERE blood_type = $2
		RETURNING id, blood_type, level, updated_at
	`, level, bloodType).Scan(&inv.ID, &inv.BloodType, &inv.Level, &inv.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("inventory not found for blood type: %s", bloodType)
		}
		return nil, err
	}
	return &inv, nil
}

// GetByBloodType returns the inventory for the specified blood type.
func (r *PgInventoryRepository) GetByBloodType(ctx context.Context, bt model.BloodType) (*model.BloodInventory, error) {
	var inv model.BloodInventory
	err := r.db.QueryRowContext(ctx, `
		SELECT id, blood_type, level, updated_at
		FROM blood_inventory
		WHERE blood_type = $1
	`, bt).Scan(&inv.ID, &inv.BloodType, &inv.Level, &inv.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// Compile-time interface check.
var _ InventoryRepo = (*PgInventoryRepository)(nil)
