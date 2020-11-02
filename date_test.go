package datatypes_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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

func TestGobEncoding(t *testing.T) {
	date := datatypes.Date(time.Now())
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(date); err != nil {
		t.Fatalf("failed to encode datatypes.Date: %v", err)
	}

	dec := gob.NewDecoder(&buf)
	var got datatypes.Date
	if err := dec.Decode(&got); err != nil {
		t.Fatalf("failed to decode to datatypes.Date: %v", err)
	}

	AssertEqual(t, date, got)
}

func TestJSONEncoding(t *testing.T) {
	date := datatypes.Date(time.Now())
	b, err := json.Marshal(date)
	if err != nil {
		t.Fatalf("failed to encode datatypes.Date: %v", err)
	}

	var got datatypes.Date
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("failed to decode to datatypes.Date: %v", err)
	}

	AssertEqual(t, date, got)
}
