package datatypes

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Date time.Time

// Scan implements sql.Scanner. Values returned from the database are always in UTC,
// as Date.Value() stores dates as timezone-free "YYYY-MM-DD" strings.
func (date *Date) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case time.Time:
		*date = Date(v)
	case string:
		t, err := time.ParseInLocation("2006-01-02", v, time.UTC)
		if err != nil {
			return err
		}
		*date = Date(t)
	case []byte:
		t, err := time.ParseInLocation("2006-01-02", string(v), time.UTC)
		if err != nil {
			return err
		}
		*date = Date(t)
	default:
		nullTime := &sql.NullTime{}
		err = nullTime.Scan(value)
		*date = Date(nullTime.Time)
	}
	return
}

// Value implements driver.Valuer. Returns the date as a "YYYY-MM-DD" string
// to prevent database drivers from applying timezone conversions that shift the date.
func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return fmt.Sprintf("%04d-%02d-%02d", y, m, d), nil
}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}

// GormDBDataType gorm db data type
func (Date) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "DATE"
	case "postgres":
		return "DATE"
	case "sqlserver":
		return "DATE"
	case "sqlite":
		return "date"
	default:
		return ""
	}
}

func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}

func (date Date) MarshalJSON() ([]byte, error) {
	return time.Time(date).MarshalJSON()
}

func (date *Date) UnmarshalJSON(b []byte) error {
	return (*time.Time)(date).UnmarshalJSON(b)
}
