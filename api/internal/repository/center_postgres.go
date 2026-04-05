package repository

import (
	"context"
	"database/sql"
	"sort"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

// PgCenterRepository is the PostgreSQL implementation of CenterRepo.
type PgCenterRepository struct {
	db *sql.DB
}

// NewPgCenterRepository creates a new PostgreSQL-backed center repository.
func NewPgCenterRepository(db *sql.DB) *PgCenterRepository {
	return &PgCenterRepository{db: db}
}

// List returns all donation centers.
func (r *PgCenterRepository) List() []*model.DonationCenter {
	ctx := context.Background()
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, address, lat, lng, capacity, available_slots, status, created_at, updated_at
		FROM donation_centers
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil
	}
	defer func() { _ = rows.Close() }()

	return scanCenters(rows)
}

// ListByDistance returns centers within radiusKm of the given coordinates, sorted by distance.
func (r *PgCenterRepository) ListByDistance(lat, lng, radiusKm float64) []*model.DonationCenter {
	all := r.List()

	type centerDist struct {
		center *model.DonationCenter
		dist   float64
	}

	candidates := make([]centerDist, 0, len(all))
	for _, c := range all {
		d := haversineKm(lat, lng, c.Lat, c.Lng)
		if d <= radiusKm {
			candidates = append(candidates, centerDist{center: c, dist: d})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].dist < candidates[j].dist
	})

	result := make([]*model.DonationCenter, len(candidates))
	for i, cd := range candidates {
		result[i] = cd.center
	}
	return result
}

// GetByID returns a donation center by ID.
func (r *PgCenterRepository) GetByID(id string) (*model.DonationCenter, error) {
	ctx := context.Background()
	var c model.DonationCenter
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, address, lat, lng, capacity, available_slots, status, created_at, updated_at
		FROM donation_centers
		WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Address, &c.Lat, &c.Lng, &c.Capacity, &c.AvailableSlots, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Create inserts a new donation center.
func (r *PgCenterRepository) Create(c *model.DonationCenter) (*model.DonationCenter, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if c.Status == "" {
		c.Status = model.CenterStatusOpen
	}

	ctx := context.Background()
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO donation_centers (name, address, lat, lng, capacity, available_slots, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`, c.Name, c.Address, c.Lat, c.Lng, c.Capacity, c.AvailableSlots, c.Status).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func scanCenters(rows *sql.Rows) []*model.DonationCenter {
	var centers []*model.DonationCenter
	for rows.Next() {
		var c model.DonationCenter
		if err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.Lat, &c.Lng, &c.Capacity, &c.AvailableSlots, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			continue
		}
		centers = append(centers, &c)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	if centers == nil {
		return []*model.DonationCenter{}
	}
	return centers
}

// Compile-time interface check.
var _ CenterRepo = (*PgCenterRepository)(nil)
