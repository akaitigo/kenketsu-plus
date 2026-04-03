package model

import "time"

type InventoryLevel string

const (
	InventoryLevelCritical   InventoryLevel = "critical"
	InventoryLevelLow        InventoryLevel = "low"
	InventoryLevelNormal     InventoryLevel = "normal"
	InventoryLevelSufficient InventoryLevel = "sufficient"
)

type BloodInventory struct {
	ID        string         `json:"id"`
	BloodType BloodType      `json:"bloodType"`
	Level     InventoryLevel `json:"level"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func IsValidInventoryLevel(level InventoryLevel) bool {
	switch level {
	case InventoryLevelCritical, InventoryLevelLow, InventoryLevelNormal, InventoryLevelSufficient:
		return true
	default:
		return false
	}
}
