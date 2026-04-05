package repository

import "github.com/akaitigo/kenketsu-plus/api/internal/model"

// CenterRepo defines the interface for donation center persistence.
type CenterRepo interface {
	List() []*model.DonationCenter
	ListByDistance(lat, lng, radiusKm float64) []*model.DonationCenter
	GetByID(id string) (*model.DonationCenter, error)
	Create(c *model.DonationCenter) (*model.DonationCenter, error)
}

// DonationRepo defines the interface for donation record persistence.
type DonationRepo interface {
	List() []*model.Donation
	Create(d *model.Donation) (*model.Donation, error)
}

// InventoryRepo defines the interface for blood inventory persistence.
type InventoryRepo interface {
	List() []*model.BloodInventory
	Update(bloodType model.BloodType, level model.InventoryLevel) (*model.BloodInventory, error)
	GetByBloodType(bt model.BloodType) (*model.BloodInventory, error)
}

// SubscriptionRepo defines the interface for push subscription persistence.
type SubscriptionRepo interface {
	List() []*model.PushSubscription
	Create(s *model.PushSubscription) (*model.PushSubscription, error)
	Delete(id string) error
	DeleteByEndpoint(endpoint string) error
}
