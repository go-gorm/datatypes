package datatypes

import (
	"database/sql/driver"
	"fmt"
	"net"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Inet implements driver.Valuer, sql.Scanner, migrator.GormDataTypeInterface and schema.GormDataType interfaces.
type Inet struct {
	net.IP
}

// Value returns value as a string.
func (ip Inet) Value() (driver.Value, error) {
	return ip.String(), nil
}

// Scan scans a string value into Inet.
func (ip *Inet) Scan(value interface{}) error {

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("can' scan: %s", value)
	}
	ip.IP = net.ParseIP(s)

	return nil
}

// GormDataType gorm common data type.
func (Inet) GormDataType() string {
	return "inet"
}

// GormDBDataType gorm db data type.
func (Inet) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "INET"
	}
	return ""
}
