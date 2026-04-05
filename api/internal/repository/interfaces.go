// Package repository defines persistence interfaces and their implementations.
package repository

import (
	"context"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

// CenterRepo defines the interface for donation center persistence.
type CenterRepo interface {
	List(ctx context.Context) []*model.DonationCenter
	ListByDistance(ctx context.Context, lat, lng, radiusKm float64) []*model.DonationCenter
	GetByID(ctx context.Context, id string) (*model.DonationCenter, error)
	Create(ctx context.Context, c *model.DonationCenter) (*model.DonationCenter, error)
}

// DonationRepo defines the interface for donation record persistence.
type DonationRepo interface {
	List(ctx context.Context) []*model.Donation
	Create(ctx context.Context, d *model.Donation) (*model.Donation, error)
}

// InventoryRepo defines the interface for blood inventory persistence.
type InventoryRepo interface {
	List(ctx context.Context) []*model.BloodInventory
	Update(ctx context.Context, bloodType model.BloodType, level model.InventoryLevel) (*model.BloodInventory, error)
	GetByBloodType(ctx context.Context, bt model.BloodType) (*model.BloodInventory, error)
}

// SubscriptionRepo defines the interface for push subscription persistence.
type SubscriptionRepo interface {
	List(ctx context.Context) []*model.PushSubscription
	Create(ctx context.Context, s *model.PushSubscription) (*model.PushSubscription, error)
	Delete(ctx context.Context, id string) error
	DeleteByEndpoint(ctx context.Context, endpoint string) error
}
