package model

import "time"

// BloodType represents a blood type classification.
type BloodType string

// BloodType constants define the ABO-Rh blood type system.
const (
	BloodTypeAPos  BloodType = "A+"
	BloodTypeANeg  BloodType = "A-"
	BloodTypeBPos  BloodType = "B+"
	BloodTypeBNeg  BloodType = "B-"
	BloodTypeOPos  BloodType = "O+"
	BloodTypeONeg  BloodType = "O-"
	BloodTypeABPos BloodType = "AB+"
	BloodTypeABNeg BloodType = "AB-"
)

// ValidBloodTypes lists all supported blood types.
var ValidBloodTypes = []BloodType{
	BloodTypeAPos, BloodTypeANeg,
	BloodTypeBPos, BloodTypeBNeg,
	BloodTypeOPos, BloodTypeONeg,
	BloodTypeABPos, BloodTypeABNeg,
}

// IsValidBloodType checks whether the given blood type is recognized.
func IsValidBloodType(bt BloodType) bool {
	for _, v := range ValidBloodTypes {
		if v == bt {
			return true
		}
	}
	return false
}

// DonationType represents the kind of blood donation performed.
type DonationType string

// DonationType constants define the supported donation types.
const (
	DonationTypeWhole400  DonationType = "whole_400"
	DonationTypeWhole200  DonationType = "whole_200"
	DonationTypeComponent DonationType = "component"
)

// Gender represents the biological sex used for donation interval calculations.
type Gender string

// Gender constants define the supported genders for donation rules.
const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

// Donation represents a single blood donation record.
type Donation struct {
	DonatedAt    time.Time    `json:"donatedAt"`
	CreatedAt    time.Time    `json:"createdAt"`
	ID           string       `json:"id"`
	BloodType    BloodType    `json:"bloodType"`
	DonationType DonationType `json:"donationType"`
	Gender       Gender       `json:"gender"`
	Memo         string       `json:"memo"`
	VolumeMl     int          `json:"volumeMl"`
}

// Validate checks that all required fields are present and valid.
func (d *Donation) Validate() error {
	if !IsValidBloodType(d.BloodType) {
		return ErrFieldInvalid("bloodType", "invalid blood type")
	}
	switch d.DonationType {
	case DonationTypeWhole400, DonationTypeWhole200, DonationTypeComponent:
	default:
		return ErrFieldInvalid("donationType", "must be whole_400, whole_200, or component")
	}
	switch d.Gender {
	case GenderMale, GenderFemale:
	default:
		return ErrFieldInvalid("gender", "must be male or female")
	}
	if d.DonatedAt.IsZero() {
		return ErrFieldRequired("donatedAt")
	}
	// M-4: component donations may have 0 volume (platelet/plasma)
	if d.DonationType != DonationTypeComponent && d.VolumeMl <= 0 {
		return ErrFieldInvalid("volumeMl", "must be positive for whole blood donations")
	}
	if d.VolumeMl < 0 {
		return ErrFieldInvalid("volumeMl", "must be non-negative")
	}
	return nil
}
