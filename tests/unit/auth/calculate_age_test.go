package auth_test

import (
	"testing"
	"time"

	"erp-service/iam/auth"

	"github.com/stretchr/testify/assert"
)

func TestCalculateAge(t *testing.T) {

	tests := []struct {
		name      string
		birthDate time.Time
		now       time.Time
		wantAge   int
	}{
		{
			name:      "standard case - 30 years old",
			birthDate: date(1994, 6, 15),
			now:       date(2024, 6, 15),
			wantAge:   30,
		},
		{
			name:      "birthday today - age equals full years",
			birthDate: date(1990, 3, 1),
			now:       date(2024, 3, 1),
			wantAge:   34,
		},
		{
			name:      "day before birthday - not yet reached this year",
			birthDate: date(1990, 3, 2),
			now:       date(2024, 3, 1),
			wantAge:   33,
		},
		{
			name:      "day after birthday - already passed this year",
			birthDate: date(1990, 2, 28),
			now:       date(2024, 3, 1),
			wantAge:   34,
		},

		{
			name: "leap-year birthday (Feb 29) - non-leap check year, before Mar 1",

			birthDate: date(2000, 2, 29),
			now:       date(2023, 2, 28),
			wantAge:   22,
		},
		{
			name: "leap-year birthday (Feb 29) - non-leap check year, on Mar 1",

			birthDate: date(2000, 2, 29),
			now:       date(2023, 3, 1),
			wantAge:   23,
		},
		{
			name:      "leap-year birthday (Feb 29) - in leap year on exact birthday",
			birthDate: date(2000, 2, 29),
			now:       date(2024, 2, 29),
			wantAge:   24,
		},
		{
			name:      "leap-year birthday (Feb 29) - in leap year, one day before",
			birthDate: date(2000, 2, 29),
			now:       date(2024, 2, 28),
			wantAge:   23,
		},
		{
			name:      "just turned 18 today - should pass minimum age check",
			birthDate: date(2006, 1, 15),
			now:       date(2024, 1, 15),
			wantAge:   18,
		},
		{
			name:      "17 years and 364 days - should fail minimum age check",
			birthDate: date(2006, 1, 16),
			now:       date(2024, 1, 15),
			wantAge:   17,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := auth.CalculateAgeAt(tt.birthDate, tt.now)
			assert.Equal(t, tt.wantAge, got,
				"CalculateAgeAt(%v, %v) = %d, want %d",
				tt.birthDate.Format("2006-01-02"),
				tt.now.Format("2006-01-02"),
				got,
				tt.wantAge,
			)
		})
	}
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
