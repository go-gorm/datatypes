package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// Scan scan value into Jsonb
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (JSON) GormDataType() string {
	return "json"
}

func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "text"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

type JSONQueryExpression struct {
	column  string
	keys    []string
	hasKeys []string
}

func JSONQuery(column string) *JSONQueryExpression {
	return &JSONQueryExpression{column: column}
}

func (jsonQuery *JSONQueryExpression) HasKey(keys ...string) *JSONQueryExpression {
	jsonQuery.hasKeys = keys
	return jsonQuery
}

func (jsonQuery *JSONQueryExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		switch stmt.Dialector.Name() {
		case "sqlite":
			// TODO
		case "mysql":
			if len(jsonQuery.hasKeys) > 0 {
				builder.WriteString(fmt.Sprintf("JSON_EXTRACT(%s, '$.%s') IS NOT NULL", stmt.Quote(jsonQuery.column), strings.Join(jsonQuery.hasKeys, ".")))
			}
		case "postgres":
			if len(jsonQuery.hasKeys) > 0 {
				stmt.WriteQuoted(jsonQuery.column)
				stmt.WriteString("::jsonb")
				for _, key := range jsonQuery.hasKeys[0 : len(jsonQuery.hasKeys)-1] {
					stmt.WriteString(" -> ")
					stmt.AddVar(stmt, key)
				}

				stmt.WriteString(" ? ")
				stmt.AddVar(stmt, jsonQuery.hasKeys[len(jsonQuery.hasKeys)-1])
			}
		}
	}
}
