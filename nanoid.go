package datatypes

import (
	"context"
	"database/sql/driver"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NanoID string

func GenerateNanoID() (string, error) {
	return gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 12)
}

func (id NanoID) GormDataType() string {
	return "VARCHAR(12)"
}

// Value implements the driver.Valuer interface.
func (id NanoID) Value() (driver.Value, error) {
	if id == "" {
		return GenerateNanoID()
	}
	idValue := string(id) // Cannot pass by ref
	return idValue, nil
}

func (id *NanoID) Scan(v interface{}) error {
	switch v := v.(type) {
	case int:
		*id = NanoID(rune(v))
		break
	case string:
		*id = NanoID(v)
	default:
		return fmt.Errorf("Could not scan: %T", v)
	}
	return nil
}

func (id NanoID) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return gorm.Expr("?", string(id))
}
