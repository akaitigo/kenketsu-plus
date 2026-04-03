package model_test

import (
	"testing"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

func TestDonationCenter_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		center  model.DonationCenter
		wantErr bool
	}{
		{
			name: "valid center",
			center: model.DonationCenter{
				Name: "東京献血ルーム", Address: "東京都千代田区",
				Lat: 35.6812, Lng: 139.7671, Capacity: 20, AvailableSlots: 5,
			},
			wantErr: false,
		},
		{name: "empty name", center: model.DonationCenter{Address: "addr", Lat: 35.0, Lng: 139.0}, wantErr: true},
		{name: "empty address", center: model.DonationCenter{Name: "name", Lat: 35.0, Lng: 139.0}, wantErr: true},
		{name: "invalid lat", center: model.DonationCenter{Name: "n", Address: "a", Lat: 91, Lng: 0}, wantErr: true},
		{name: "invalid lng", center: model.DonationCenter{Name: "n", Address: "a", Lat: 0, Lng: 181}, wantErr: true},
		{name: "negative capacity", center: model.DonationCenter{Name: "n", Address: "a", Lat: 0, Lng: 0, Capacity: -1}, wantErr: true},
		{
			name:    "slots exceed capacity",
			center:  model.DonationCenter{Name: "n", Address: "a", Lat: 0, Lng: 0, Capacity: 5, AvailableSlots: 10},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.center.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDonation_Validate(t *testing.T) {
	t.Parallel()

	valid := model.Donation{
		BloodType: model.BloodTypeAPos, DonationType: model.DonationTypeWhole400,
		Gender: model.GenderMale, DonatedAt: time.Now(), VolumeMl: 400,
	}

	tests := []struct {
		modify  func(d *model.Donation)
		name    string
		wantErr bool
	}{
		{name: "valid", modify: func(_ *model.Donation) {}, wantErr: false},
		{name: "invalid blood type", modify: func(d *model.Donation) { d.BloodType = "X" }, wantErr: true},
		{name: "invalid donation type", modify: func(d *model.Donation) { d.DonationType = "bad" }, wantErr: true},
		{name: "invalid gender", modify: func(d *model.Donation) { d.Gender = "other" }, wantErr: true},
		{name: "zero donated_at", modify: func(d *model.Donation) { d.DonatedAt = time.Time{} }, wantErr: true},
		{name: "zero volume whole blood", modify: func(d *model.Donation) { d.VolumeMl = 0 }, wantErr: true},
		{name: "negative volume", modify: func(d *model.Donation) { d.VolumeMl = -1 }, wantErr: true},
		{name: "component zero volume ok", modify: func(d *model.Donation) { d.DonationType = model.DonationTypeComponent; d.VolumeMl = 0 }, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := valid
			tt.modify(&d)
			err := d.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPushSubscription_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		sub     model.PushSubscription
		name    string
		wantErr bool
	}{
		{name: "valid", sub: model.PushSubscription{Endpoint: "https://example.com", P256dh: "key", Auth: "auth"}, wantErr: false},
		{name: "empty endpoint", sub: model.PushSubscription{P256dh: "key", Auth: "auth"}, wantErr: true},
		{name: "empty p256dh", sub: model.PushSubscription{Endpoint: "https://example.com", Auth: "auth"}, wantErr: true},
		{name: "empty auth", sub: model.PushSubscription{Endpoint: "https://example.com", P256dh: "key"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.sub.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidBloodType(t *testing.T) {
	t.Parallel()

	if !model.IsValidBloodType(model.BloodTypeAPos) {
		t.Error("expected A+ to be valid")
	}
	if model.IsValidBloodType("X+") {
		t.Error("expected X+ to be invalid")
	}
}

func TestIsValidInventoryLevel(t *testing.T) {
	t.Parallel()

	if !model.IsValidInventoryLevel(model.InventoryLevelCritical) {
		t.Error("expected critical to be valid")
	}
	if model.IsValidInventoryLevel("unknown") {
		t.Error("expected unknown to be invalid")
	}
}
