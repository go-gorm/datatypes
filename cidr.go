package datatypes

import (
	"database/sql/driver"
	"fmt"
	"net"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// CIDR implements driver.Valuer, sql.Scanner, migrator.GormDataTypeInterface and schema.GormDataType interfaces.
type CIDR struct {
	net.IPNet
}

// Value returns value as a string.
func (n CIDR) Value() (driver.Value, error) {
	return n.String(), nil
}

// Scan scans a string value into CIDR.
func (n *CIDR) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("can't scan: %v", value)
	}

	_, cidr, err := net.ParseCIDR(s)
	if err != nil {
		return fmt.Errorf("can't parse cidr %q: %w", s, err)
	}
	n.IPNet = *cidr
	return nil
}

// GormDataType gorm common data type.
func (CIDR) GormDataType() string {
	return "cidr"
}

// GormDBDataType gorm db data type.
func (CIDR) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "CIDR"
	}
	return ""
}
