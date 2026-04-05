package repository

import (
	"context"
	"database/sql"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

// PgDonationRepository is the PostgreSQL implementation of DonationRepo.
type PgDonationRepository struct {
	db *sql.DB
}

// NewPgDonationRepository creates a new PostgreSQL-backed donation repository.
func NewPgDonationRepository(db *sql.DB) *PgDonationRepository {
	return &PgDonationRepository{db: db}
}

// List returns all donations ordered by donated_at desc.
func (r *PgDonationRepository) List(ctx context.Context) []*model.Donation {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, blood_type, donation_type, gender, donated_at, volume_ml, COALESCE(memo, ''), created_at
		FROM donations
		ORDER BY donated_at DESC
	`)
	if err != nil {
		return nil
	}
	defer func() { _ = rows.Close() }()

	var donations []*model.Donation
	for rows.Next() {
		var d model.Donation
		if err := rows.Scan(&d.ID, &d.BloodType, &d.DonationType, &d.Gender, &d.DonatedAt, &d.VolumeMl, &d.Memo, &d.CreatedAt); err != nil {
			continue
		}
		donations = append(donations, &d)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	if donations == nil {
		return []*model.Donation{}
	}
	return donations
}

// Create inserts a new donation record.
func (r *PgDonationRepository) Create(ctx context.Context, d *model.Donation) (*model.Donation, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO donations (blood_type, donation_type, gender, donated_at, volume_ml, memo)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`, d.BloodType, d.DonationType, d.Gender, d.DonatedAt, d.VolumeMl, d.Memo).Scan(&d.ID, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Compile-time interface check.
var _ DonationRepo = (*PgDonationRepository)(nil)
