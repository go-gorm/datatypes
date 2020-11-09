package datatypes

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// Date defiend Date data type, need to implements driver.Valuer, sql.Scanner interface
type Date struct {
	time.Time
}

// Value return date value, implement driver.Valuer interface
func (date Date) Value() (driver.Value, error) {
	y, m, d := date.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, date.Location()), nil
}

// Scan scan value into time.Time, implements sql.Scanner interface
func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date{nullTime.Time}
	return
}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}
