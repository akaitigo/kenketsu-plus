package model

import "time"

// InventoryLevel represents the urgency level of blood inventory.
type InventoryLevel string

// InventoryLevel constants define the possible inventory states.
const (
	InventoryLevelCritical   InventoryLevel = "critical"
	InventoryLevelLow        InventoryLevel = "low"
	InventoryLevelNormal     InventoryLevel = "normal"
	InventoryLevelSufficient InventoryLevel = "sufficient"
)

// BloodInventory represents the current stock level of a specific blood type.
type BloodInventory struct {
	UpdatedAt time.Time      `json:"updatedAt"`
	ID        string         `json:"id"`
	BloodType BloodType      `json:"bloodType"`
	Level     InventoryLevel `json:"level"`
}

// IsValidInventoryLevel checks whether the given inventory level is recognized.
func IsValidInventoryLevel(level InventoryLevel) bool {
	switch level {
	case InventoryLevelCritical, InventoryLevelLow, InventoryLevelNormal, InventoryLevelSufficient:
		return true
	default:
		return false
	}
}
