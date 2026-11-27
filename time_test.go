package datatypes_test

import (
	"testing"

	"gorm.io/datatypes"
	. "gorm.io/gorm/utils/tests"
)

func TestTime(t *testing.T) {
	if SupportedDriver("mysql", "postgres", "sqlite", "sqlserver") {
		type UserWithTime struct {
			ID   uint
			Name string
			Time datatypes.Time
		}

		DB.Migrator().DropTable(&UserWithTime{})
		if err := DB.Migrator().AutoMigrate(&UserWithTime{}); err != nil {
			t.Fatalf("failed to migrate, got error: %v", err)
		}

		user := UserWithTime{Name: "user1", Time: datatypes.NewTime(1, 2, 3, 0)}
		DB.Create(&user)

		result := UserWithTime{}
		if err := DB.First(&result, "name = ? AND time = ?", "user1", datatypes.NewTime(1, 2, 3, 0)).Error; err != nil {
			t.Fatalf("failed to find record with time, got error: %v", err)
		}

		AssertEqual(t, result.Time, datatypes.NewTime(1, 2, 3, 0))
	}
}

func TestScanFractionalSeconds(t *testing.T) {
	for _, tc := range []struct {
		input string
		nsec  int
	}{
		{"12:34:56", 0},
		{"12:34:56.1", 100000000},
		{"12:34:56.004002", 4002000},
		{"12:34:56.123456789", 123456789},
	} {
		want := datatypes.NewTime(12, 34, 56, tc.nsec)
		for _, input := range []interface{}{tc.input, []byte(tc.input)} {
			var got datatypes.Time
			if err := got.Scan(input); err != nil || got != want {
				t.Errorf("Scan(%v) = %v, %v; want %v", input, got, err, want)
			}
		}
	}
}
