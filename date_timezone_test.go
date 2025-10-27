package datatypes_test

import (
	"testing"
	"time"

	"gorm.io/datatypes"
)

func TestDateTimezoneHandling(t *testing.T) {
	// Test that Date.Value() returns consistent results regardless of timezone
	// This addresses the issue with PreferSimpleProtocol in PostgreSQL

	// Create a time in a timezone that's ahead of UTC
	loc, err := time.LoadLocation("Europe/Berlin") // UTC+1 or UTC+2
	if err != nil {
		t.Skip("Could not load Europe/Berlin timezone")
	}

	// Create a time that would be the previous day in UTC
	// For example: 2025-09-18 01:00:00 +02:00 is 2025-09-17 23:00:00 UTC
	localTime := time.Date(2025, 9, 18, 1, 0, 0, 0, loc)
	date := datatypes.Date(localTime)

	// Get the Value() result
	value, err := date.Value()
	if err != nil {
		t.Fatalf("date.Value() returned error: %v", err)
	}

	// The value should be a time.Time
	timeValue, ok := value.(time.Time)
	if !ok {
		t.Fatalf("date.Value() should return time.Time, got %T", value)
	}

	// The date part should match the original date (2025-09-18)
	// regardless of timezone, and should be in UTC
	expectedYear, expectedMonth, expectedDay := 2025, time.September, 18
	actualYear, actualMonth, actualDay := timeValue.Date()

	if actualYear != expectedYear || actualMonth != expectedMonth || actualDay != expectedDay {
		t.Errorf("Expected date %d-%02d-%02d, got %d-%02d-%02d",
			expectedYear, expectedMonth, expectedDay,
			actualYear, actualMonth, actualDay)
	}

	// The time should be in UTC
	if timeValue.Location() != time.UTC {
		t.Errorf("Expected UTC timezone, got %v", timeValue.Location())
	}

	// The time should be midnight (00:00:00)
	hour, min, sec := timeValue.Clock()
	if hour != 0 || min != 0 || sec != 0 {
		t.Errorf("Expected midnight (00:00:00), got %02d:%02d:%02d", hour, min, sec)
	}
}

func TestDateValueConsistency(t *testing.T) {
	// Test that the same date in different timezones produces the same Value()
	date := time.Date(2025, 9, 18, 15, 30, 45, 0, time.UTC)

	// Create the same date in different timezones
	utcDate := datatypes.Date(date)
	
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("Could not load America/New_York timezone")
	}
	nyDate := datatypes.Date(date.In(loc))

	// Both should produce the same Value()
	utcValue, err := utcDate.Value()
	if err != nil {
		t.Fatalf("utcDate.Value() returned error: %v", err)
	}

	nyValue, err := nyDate.Value()
	if err != nil {
		t.Fatalf("nyDate.Value() returned error: %v", err)
	}

	// Both values should be equal
	if !utcValue.(time.Time).Equal(nyValue.(time.Time)) {
		t.Errorf("Expected same value for same date in different timezones, got %v != %v", utcValue, nyValue)
	}
}