// Package model defines domain types and validation rules for kenketsu-plus.
package model

import "time"

// CenterStatus represents the operational status of a donation center.
type CenterStatus string

// CenterStatus constants define the possible statuses for a donation center.
const (
	CenterStatusOpen   CenterStatus = "open"
	CenterStatusClosed CenterStatus = "closed"
	CenterStatusFull   CenterStatus = "full"
)

// DonationCenter represents a blood donation facility with location and availability.
type DonationCenter struct {
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Address        string       `json:"address"`
	Status         CenterStatus `json:"status"`
	Lat            float64      `json:"lat"`
	Lng            float64      `json:"lng"`
	Capacity       int          `json:"capacity"`
	AvailableSlots int          `json:"availableSlots"`
}

// Validate checks that all required fields are present and within valid ranges.
func (c *DonationCenter) Validate() error {
	if c.Name == "" {
		return ErrFieldRequired("name")
	}
	if c.Address == "" {
		return ErrFieldRequired("address")
	}
	if c.Lat < -90 || c.Lat > 90 {
		return ErrFieldInvalid("lat", "must be between -90 and 90")
	}
	if c.Lng < -180 || c.Lng > 180 {
		return ErrFieldInvalid("lng", "must be between -180 and 180")
	}
	if c.Capacity < 0 {
		return ErrFieldInvalid("capacity", "must be non-negative")
	}
	if c.AvailableSlots < 0 || c.AvailableSlots > c.Capacity {
		return ErrFieldInvalid("availableSlots", "must be between 0 and capacity")
	}
	return nil
}
