package datatypes

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// JSONType give a generic data type for json encoded data.
type JSONType[T any] struct {
	Data T
}

// Value return json value, implement driver.Valuer interface
func (j JSONType[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

// Scan scan value into JSONType[T], implements sql.Scanner interface
func (j *JSONType[T]) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(bytes, &j.Data)
}

// MarshalJSON to output non base64 encoded []byte
func (j JSONType[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Data)
}

// UnmarshalJSON to deserialize []byte
func (j *JSONType[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &j.Data)
}

// GormDataType gorm common data type
func (JSONType[T]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (JSONType[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (js JSONType[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := js.MarshalJSON()

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}
