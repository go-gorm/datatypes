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

// JSON defined JSON data type, need to implements driver.Valuer, sql.Scanner interface
type JSON json.RawMessage

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	bytes, err := json.RawMessage(j).MarshalJSON()
	return string(bytes), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = JSON("null")
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (j JSON) MarshalJSON() ([]byte, error) {
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON to deserialize []byte
func (j *JSON) UnmarshalJSON(b []byte) error {
	result := json.RawMessage{}
	err := result.UnmarshalJSON(b)
	*j = JSON(result)
	return err
}

func (j JSON) String() string {
	return string(j)
}

// GormDataType gorm common data type
func (JSON) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (js JSON) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(js) == 0 {
		return gorm.Expr("NULL")
	}

	data, _ := js.MarshalJSON()

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

// JSONQueryExpression json query expression, implements clause.Expression interface to use as querier
type JSONQueryExpression struct {
	column      string
	keys        []string
	hasKeys     bool
	equals      bool
	equalsValue interface{}
}

// JSONQuery query column as json
func JSONQuery(column string) *JSONQueryExpression {
	return &JSONQueryExpression{column: column}
}

// HasKey returns clause.Expression
func (jsonQuery *JSONQueryExpression) HasKey(keys ...string) *JSONQueryExpression {
	jsonQuery.keys = keys
	jsonQuery.hasKeys = true
	return jsonQuery
}

// Keys returns clause.Expression
func (jsonQuery *JSONQueryExpression) Equals(value interface{}, keys ...string) *JSONQueryExpression {
	jsonQuery.keys = keys
	jsonQuery.equals = true
	jsonQuery.equalsValue = value
	return jsonQuery
}

// Build implements clause.Expression
func (jsonQuery *JSONQueryExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		switch stmt.Dialector.Name() {
		case "mysql", "sqlite":
			switch {
			case jsonQuery.hasKeys:
				if len(jsonQuery.keys) > 0 {
					builder.WriteString("JSON_EXTRACT(" + stmt.Quote(jsonQuery.column) + ",")
					builder.AddVar(stmt, "$."+strings.Join(jsonQuery.keys, "."))
					builder.WriteString(") IS NOT NULL")
				}
			case jsonQuery.equals:
				if len(jsonQuery.keys) > 0 {
					builder.WriteString("JSON_EXTRACT(" + stmt.Quote(jsonQuery.column) + ",")
					builder.AddVar(stmt, "$."+strings.Join(jsonQuery.keys, "."))
					builder.WriteString(") = ")
					if _, ok := jsonQuery.equalsValue.(bool); ok {
						builder.WriteString(fmt.Sprint(jsonQuery.equalsValue))
					} else {
						stmt.AddVar(builder, jsonQuery.equalsValue)
					}
				}
			}
		case "postgres":
			switch {
			case jsonQuery.hasKeys:
				if len(jsonQuery.keys) > 0 {
					stmt.WriteQuoted(jsonQuery.column)
					stmt.WriteString("::jsonb")
					for _, key := range jsonQuery.keys[0 : len(jsonQuery.keys)-1] {
						stmt.WriteString(" -> ")
						stmt.AddVar(builder, key)
					}

					stmt.WriteString(" ? ")
					stmt.AddVar(builder, jsonQuery.keys[len(jsonQuery.keys)-1])
				}
			case jsonQuery.equals:
				if len(jsonQuery.keys) > 0 {
					builder.WriteString(fmt.Sprintf("json_extract_path_text(%v::json,", stmt.Quote(jsonQuery.column)))

					for idx, key := range jsonQuery.keys {
						if idx > 0 {
							builder.WriteByte(',')
						}
						stmt.AddVar(builder, key)
					}
					builder.WriteString(") = ")

					if _, ok := jsonQuery.equalsValue.(string); ok {
						stmt.AddVar(builder, jsonQuery.equalsValue)
					} else {
						stmt.AddVar(builder, fmt.Sprint(jsonQuery.equalsValue))
					}
				}
			}
		}
	}
}
