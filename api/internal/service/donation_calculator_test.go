package service_test

import (
	"testing"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/service"
)

func makeDonation(donationType model.DonationType, gender model.Gender, donatedAt time.Time) *model.Donation {
	return &model.Donation{
		DonatedAt:    donatedAt,
		BloodType:    model.BloodTypeAPos,
		DonationType: donationType,
		Gender:       gender,
		VolumeMl:     400,
	}
}

func TestNextAvailableDate_NoDonations(t *testing.T) {
	t.Parallel()
	calc := service.NewDonationCalculator()
	result := calc.NextAvailableDate(nil, model.GenderMale)
	if !result.CanDonateToday {
		t.Error("expected can donate today with no records")
	}
	if result.DaysRemaining != 0 {
		t.Errorf("expected 0 days remaining, got %d", result.DaysRemaining)
	}
}

func TestNextAvailableDate_Whole400_Male(t *testing.T) {
	t.Parallel()
	calc := service.NewDonationCalculator()
	donations := []*model.Donation{
		makeDonation(model.DonationTypeWhole400, model.GenderMale, time.Now().AddDate(0, 0, -7)),
	}
	result := calc.NextAvailableDate(donations, model.GenderMale)
	if result.CanDonateToday {
		t.Error("should not be able to donate 7 days after whole 400ml (male needs 12 weeks)")
	}
	expectedDays := 12*7 - 7
	if result.DaysRemaining < expectedDays-1 || result.DaysRemaining > expectedDays+1 {
		t.Errorf("expected ~%d days remaining, got %d", expectedDays, result.DaysRemaining)
	}
}

func TestNextAvailableDate_Whole400_Female(t *testing.T) {
	t.Parallel()
	calc := service.NewDonationCalculator()
	donations := []*model.Donation{
		makeDonation(model.DonationTypeWhole400, model.GenderFemale, time.Now().AddDate(0, 0, -7)),
	}
	result := calc.NextAvailableDate(donations, model.GenderFemale)
	if result.CanDonateToday {
		t.Error("should not be able to donate 7 days after whole 400ml (female needs 16 weeks)")
	}
	expectedDays := 16*7 - 7
	if result.DaysRemaining < expectedDays-1 || result.DaysRemaining > expectedDays+1 {
		t.Errorf("expected ~%d days remaining, got %d", expectedDays, result.DaysRemaining)
	}
}

func TestNextAvailableDate_Whole200(t *testing.T) {
	t.Parallel()
	calc := service.NewDonationCalculator()
	donations := []*model.Donation{
		makeDonation(model.DonationTypeWhole200, model.GenderMale, time.Now().AddDate(0, 0, -7)),
	}
	result := calc.NextAvailableDate(donations, model.GenderMale)
	if result.CanDonateToday {
		t.Error("should not be able to donate 7 days after whole 200ml (needs 4 weeks)")
	}
	expectedDays := 4*7 - 7
	if result.DaysRemaining < expectedDays-1 || result.DaysRemaining > expectedDays+1 {
		t.Errorf("expected ~%d days remaining, got %d", expectedDays, result.DaysRemaining)
	}
}

func TestNextAvailableDate_Component(t *testing.T) {
	t.Parallel()
	calc := service.NewDonationCalculator()
	donations := []*model.Donation{
		makeDonation(model.DonationTypeComponent, model.GenderMale, time.Now().AddDate(0, 0, -7)),
	}
	result := calc.NextAvailableDate(donations, model.GenderMale)
	if result.CanDonateToday {
		t.Error("should not be able to donate 7 days after component (needs 2 weeks)")
	}
	expectedDays := 2*7 - 7
	if result.DaysRemaining < expectedDays-1 || result.DaysRemaining > expectedDays+1 {
		t.Errorf("expected ~%d days remaining, got %d", expectedDays, result.DaysRemaining)
	}
}

func TestNextAvailableDate_IntervalExpired(t *testing.T) {
	t.Parallel()
	calc := service.NewDonationCalculator()
	donations := []*model.Donation{
		makeDonation(model.DonationTypeWhole400, model.GenderMale, time.Now().AddDate(0, 0, -100)),
	}
	result := calc.NextAvailableDate(donations, model.GenderMale)
	if !result.CanDonateToday {
		t.Error("should be able to donate after 100 days (12 weeks = 84 days)")
	}
}

func TestNextAvailableDate_AnnualLimit_Male(t *testing.T) {
	t.Parallel()
	calc := service.NewDonationCalculator()
	now := time.Now()
	donations := []*model.Donation{
		makeDonation(model.DonationTypeWhole400, model.GenderMale, now.AddDate(0, -10, 0)),
		makeDonation(model.DonationTypeWhole400, model.GenderMale, now.AddDate(0, -7, 0)),
		makeDonation(model.DonationTypeWhole400, model.GenderMale, now.AddDate(0, -4, 0)),
	}
	result := calc.NextAvailableDate(donations, model.GenderMale)
	if result.CanDonateToday {
		t.Error("male should not be able to donate with 3 whole donations in the year")
	}
}

func TestNextAvailableDate_AnnualLimit_Female(t *testing.T) {
	t.Parallel()
	calc := service.NewDonationCalculator()
	now := time.Now()
	donations := []*model.Donation{
		makeDonation(model.DonationTypeWhole400, model.GenderFemale, now.AddDate(0, -8, 0)),
		makeDonation(model.DonationTypeWhole400, model.GenderFemale, now.AddDate(0, -4, 0)),
	}
	result := calc.NextAvailableDate(donations, model.GenderFemale)
	if result.CanDonateToday {
		t.Error("female should not be able to donate with 2 whole donations in the year")
	}
}

// TestNextAvailableDate_EdgeCases exercises mixed donation-type sequences,
// annual-limit interactions, boundary values, unsorted input, and unknown types
// via a table (#23). Day windows are inclusive and allow ±2 slack to absorb
// time-of-day and timezone effects on the day-count arithmetic.
func TestNextAvailableDate_EdgeCases(t *testing.T) {
	t.Parallel()

	calc := service.NewDonationCalculator()

	d := func(dt model.DonationType, g model.Gender, daysAgo int) *model.Donation {
		return makeDonation(dt, g, time.Now().AddDate(0, 0, -daysAgo))
	}

	tests := []struct {
		name          string
		gender        model.Gender
		donations     []*model.Donation
		wantDaysLo    int // asserted only when wantCanDonate is false
		wantDaysHi    int
		wantCanDonate bool
	}{
		{
			name: "component after whole400 uses 2-week interval (component is last)",
			donations: []*model.Donation{
				d(model.DonationTypeWhole400, model.GenderMale, 20),
				d(model.DonationTypeComponent, model.GenderMale, 5),
			},
			gender:        model.GenderMale,
			wantCanDonate: false,
			wantDaysLo:    7,
			wantDaysHi:    11,
		},
		{
			name: "component to whole400 to component, last component interval expired",
			donations: []*model.Donation{
				d(model.DonationTypeComponent, model.GenderMale, 60),
				d(model.DonationTypeWhole400, model.GenderMale, 40),
				d(model.DonationTypeComponent, model.GenderMale, 20),
			},
			gender:        model.GenderMale,
			wantCanDonate: true,
		},
		{
			name: "whole400 then component, component still in interval",
			donations: []*model.Donation{
				d(model.DonationTypeWhole400, model.GenderMale, 30),
				d(model.DonationTypeComponent, model.GenderMale, 3),
			},
			gender:        model.GenderMale,
			wantCanDonate: false,
			wantDaysLo:    9,
			wantDaysHi:    13,
		},
		{
			name: "male annual whole limit dominates over interval",
			donations: []*model.Donation{
				d(model.DonationTypeWhole400, model.GenderMale, 300),
				d(model.DonationTypeWhole400, model.GenderMale, 200),
				d(model.DonationTypeWhole400, model.GenderMale, 100),
			},
			gender:        model.GenderMale,
			wantCanDonate: false,
			wantDaysLo:    62,
			wantDaysHi:    68,
		},
		{
			name: "whole200 counts toward annual whole limit",
			donations: []*model.Donation{
				d(model.DonationTypeWhole200, model.GenderMale, 300),
				d(model.DonationTypeWhole400, model.GenderMale, 200),
				d(model.DonationTypeWhole400, model.GenderMale, 100),
			},
			gender:        model.GenderMale,
			wantCanDonate: false,
			wantDaysLo:    62,
			wantDaysHi:    68,
		},
		{
			name: "component donations never hit annual limit",
			donations: []*model.Donation{
				d(model.DonationTypeComponent, model.GenderMale, 100),
				d(model.DonationTypeComponent, model.GenderMale, 80),
				d(model.DonationTypeComponent, model.GenderMale, 60),
				d(model.DonationTypeComponent, model.GenderMale, 40),
				d(model.DonationTypeComponent, model.GenderMale, 20),
			},
			gender:        model.GenderMale,
			wantCanDonate: true,
		},
		{
			name: "female annual whole limit is two",
			donations: []*model.Donation{
				d(model.DonationTypeWhole400, model.GenderFemale, 200),
				d(model.DonationTypeWhole400, model.GenderFemale, 100),
			},
			gender:        model.GenderFemale,
			wantCanDonate: false,
			wantDaysLo:    160,
			wantDaysHi:    170,
		},
		{
			name: "female single whole400 uses 16-week interval",
			donations: []*model.Donation{
				d(model.DonationTypeWhole400, model.GenderFemale, 50),
			},
			gender:        model.GenderFemale,
			wantCanDonate: false,
			wantDaysLo:    59,
			wantDaysHi:    65,
		},
		{
			name: "unsorted input still uses the most recent donation",
			donations: []*model.Donation{
				d(model.DonationTypeComponent, model.GenderMale, 3),
				d(model.DonationTypeWhole400, model.GenderMale, 50),
			},
			gender:        model.GenderMale,
			wantCanDonate: false,
			wantDaysLo:    9,
			wantDaysHi:    13,
		},
		{
			name: "unknown donation type defaults to 12-week interval",
			donations: []*model.Donation{
				d(model.DonationType("unknown"), model.GenderMale, 10),
			},
			gender:        model.GenderMale,
			wantCanDonate: false,
			wantDaysLo:    71,
			wantDaysHi:    77,
		},
		{
			name: "boundary just past interval can donate",
			donations: []*model.Donation{
				d(model.DonationTypeWhole400, model.GenderMale, 88),
			},
			gender:        model.GenderMale,
			wantCanDonate: true,
		},
		{
			name: "boundary just before interval cannot donate",
			donations: []*model.Donation{
				d(model.DonationTypeWhole400, model.GenderMale, 80),
			},
			gender:        model.GenderMale,
			wantCanDonate: false,
			wantDaysLo:    2,
			wantDaysHi:    6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := calc.NextAvailableDate(tt.donations, tt.gender)

			if result.CanDonateToday != tt.wantCanDonate {
				t.Fatalf("CanDonateToday = %v, want %v (reason: %q, daysRemaining: %d)",
					result.CanDonateToday, tt.wantCanDonate, result.Reason, result.DaysRemaining)
			}

			if tt.wantCanDonate {
				if result.DaysRemaining != 0 {
					t.Errorf("DaysRemaining = %d, want 0 when donation is possible", result.DaysRemaining)
				}
				return
			}

			if result.DaysRemaining < tt.wantDaysLo || result.DaysRemaining > tt.wantDaysHi {
				t.Errorf("DaysRemaining = %d, want within [%d, %d]",
					result.DaysRemaining, tt.wantDaysLo, tt.wantDaysHi)
			}
		})
	}
}

func TestNextAvailableDate_YearBoundary(t *testing.T) {
	t.Parallel()
	calc := service.NewDonationCalculator()
	// 年末に全血400ml献血した場合、12週後は翌年になる
	// 献血日から12週 = 84日後を計算
	donationDate := time.Now().AddDate(0, 0, -14) // 2週間前
	donations := []*model.Donation{
		makeDonation(model.DonationTypeWhole400, model.GenderMale, donationDate),
	}
	result := calc.NextAvailableDate(donations, model.GenderMale)
	expected := donationDate.AddDate(0, 0, 12*7)
	if result.CanDonateToday {
		t.Error("should not be able to donate 14 days after whole 400ml")
	}
	if result.NextDate.Before(expected.AddDate(0, 0, -1)) || result.NextDate.After(expected.AddDate(0, 0, 1)) {
		t.Errorf("expected next date around %v, got %v", expected, result.NextDate)
	}
}
