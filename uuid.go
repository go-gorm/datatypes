package datatypes

import (
	"database/sql/driver"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type UUID uuid.UUID

// GormDataType gorm common data type
func (UUID) GormDataType() string {
	return "string"
}

// GormDBDataType gorm db data type
func (UUID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "LONGTEXT"
	case "postgres":
		return "UUID"
	case "sqlserver":
		return "NVARCHAR"
	case "sqlite":
		return "TEXT"
	default:
		return ""
	}
}

func (u *UUID) Scan(value interface{}) error {
	var result uuid.UUID
	if err := result.Scan(value); err != nil {
		return err
	}
	*u = UUID(result)
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u).Value()
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func (u UUID) Equals(other UUID) bool {
	return u.String() == other.String()
}

func (u UUID) Length() int {
	return len(u.String())
}

func (u UUID) IsNil() bool {
	return uuid.UUID(u) == uuid.Nil
}

func (u UUID) IsEmpty() bool {
	return u.IsNil() || u.Length() == 0
}

func (u *UUID) IsNilPtr() bool {
	return u == nil
}

func (u *UUID) IsEmptyPtr() bool {
	return u.IsNilPtr() || u.IsEmpty()
}
