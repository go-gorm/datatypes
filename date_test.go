package datatypes_test

import (
	"testing"
	"time"

	"github.com/jinzhu/now"
	"gorm.io/datatypes"
	. "gorm.io/gorm/utils/tests"
)

func TestDate(t *testing.T) {
	type UserWithDate struct {
		ID   uint
		Name string
		Date datatypes.Date
	}

	DB.Migrator().DropTable(&UserWithDate{})
	if err := DB.Migrator().AutoMigrate(&UserWithDate{}); err != nil {
		t.Errorf("failed to migrate, got error: %v", err)
	}

	curTime := time.Now().UTC()
	beginningOfDay := now.New(curTime).BeginningOfDay()

	user := UserWithDate{Name: "jinzhu", Date: datatypes.Date(curTime)}
	DB.Create(&user)

	result := UserWithDate{}
	if err := DB.First(&result, "name = ? AND date = ?", "jinzhu", datatypes.Date(curTime)).Error; err != nil {
		t.Fatalf("Failed to find record with date")
	}

	AssertEqual(t, result.Date, beginningOfDay)
}
