package model

import "time"

type BloodType string

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

var ValidBloodTypes = []BloodType{
	BloodTypeAPos, BloodTypeANeg,
	BloodTypeBPos, BloodTypeBNeg,
	BloodTypeOPos, BloodTypeONeg,
	BloodTypeABPos, BloodTypeABNeg,
}

func IsValidBloodType(bt BloodType) bool {
	for _, v := range ValidBloodTypes {
		if v == bt {
			return true
		}
	}
	return false
}

type DonationType string

const (
	DonationTypeWhole400  DonationType = "whole_400"
	DonationTypeWhole200  DonationType = "whole_200"
	DonationTypeComponent DonationType = "component"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

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
