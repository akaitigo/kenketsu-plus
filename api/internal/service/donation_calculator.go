package service

import (
	"sort"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

type DonationCalculator struct{}

func NewDonationCalculator() *DonationCalculator {
	return &DonationCalculator{}
}

type NextAvailableResult struct {
	NextDate       time.Time `json:"nextDate"`
	DaysRemaining  int       `json:"daysRemaining"`
	CanDonateToday bool      `json:"canDonateToday"`
	Reason         string    `json:"reason"`
}

func (c *DonationCalculator) intervalWeeks(donationType model.DonationType, gender model.Gender) int {
	switch donationType {
	case model.DonationTypeWhole400:
		if gender == model.GenderMale {
			return 12
		}
		return 16
	case model.DonationTypeWhole200:
		return 4
	case model.DonationTypeComponent:
		return 2
	default:
		return 12
	}
}

func (c *DonationCalculator) annualWholeLimit(gender model.Gender) int {
	if gender == model.GenderMale {
		return 3
	}
	return 2
}

func (c *DonationCalculator) NextAvailableDate(donations []*model.Donation, gender model.Gender) NextAvailableResult {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if len(donations) == 0 {
		return NextAvailableResult{
			NextDate:       today,
			DaysRemaining:  0,
			CanDonateToday: true,
			Reason:         "献血記録がありません。本日から献血可能です。",
		}
	}

	sorted := make([]*model.Donation, len(donations))
	copy(sorted, donations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].DonatedAt.After(sorted[j].DonatedAt)
	})

	last := sorted[0]
	weeks := c.intervalWeeks(last.DonationType, gender)
	intervalNext := last.DonatedAt.AddDate(0, 0, weeks*7)

	oneYearAgo := today.AddDate(-1, 0, 0)
	wholeCount := 0
	for _, d := range sorted {
		if d.DonatedAt.After(oneYearAgo) &&
			(d.DonationType == model.DonationTypeWhole400 || d.DonationType == model.DonationTypeWhole200) {
			wholeCount++
		}
	}

	limit := c.annualWholeLimit(gender)
	if wholeCount >= limit {
		oldestInYear := findOldestWholeInYear(sorted, oneYearAgo)
		if oldestInYear != nil {
			annualNext := oldestInYear.DonatedAt.AddDate(1, 0, 0)
			if annualNext.After(intervalNext) {
				return NextAvailableResult{
					NextDate:       annualNext,
					DaysRemaining:  daysUntil(today, annualNext),
					CanDonateToday: false,
					Reason:         "年間全血回数上限に達しています。",
				}
			}
		}
	}

	if intervalNext.After(today) {
		return NextAvailableResult{
			NextDate:       intervalNext,
			DaysRemaining:  daysUntil(today, intervalNext),
			CanDonateToday: false,
			Reason:         "前回の献血からの間隔制限中です。",
		}
	}

	return NextAvailableResult{
		NextDate:       today,
		DaysRemaining:  0,
		CanDonateToday: true,
		Reason:         "本日から献血可能です。",
	}
}

func findOldestWholeInYear(sorted []*model.Donation, oneYearAgo time.Time) *model.Donation {
	var oldest *model.Donation
	for _, d := range sorted {
		if d.DonatedAt.After(oneYearAgo) &&
			(d.DonationType == model.DonationTypeWhole400 || d.DonationType == model.DonationTypeWhole200) {
			oldest = d
		}
	}
	return oldest
}

func daysUntil(from, to time.Time) int {
	d := to.Sub(from).Hours() / 24
	if d < 0 {
		return 0
	}
	return int(d)
}
