package service_test

import (
	"testing"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/service"
)

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

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
