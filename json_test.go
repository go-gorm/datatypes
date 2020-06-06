package datatypes_test

import (
	"database/sql/driver"
	"fmt"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	. "gorm.io/gorm/utils/tests"
)

var _ driver.Valuer = &datatypes.JSON{}

func TestJSON(t *testing.T) {
	if !SupportedDriver("mysql", "postgres") {
		fmt.Println(DB.Dialector.Name())
		return
	}

	DB.Dialector.Name()
	type UserWithJSON struct {
		gorm.Model
		Name       string
		Attributes datatypes.JSON
	}

	DB.Migrator().DropTable(&UserWithJSON{})
	if err := DB.Migrator().AutoMigrate(&UserWithJSON{}); err != nil {
		t.Errorf("failed to migrate, got error: %v", err)
	}

	users := []UserWithJSON{{
		Name:       "json-1",
		Attributes: datatypes.JSON([]byte(`{"name": "json-1", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
	}, {
		Name:       "json-2",
		Attributes: datatypes.JSON([]byte(`{"name": "json-2", "age": 28, "tags": ["tag1", "tag3"], "role": "admin", "orgs": {"orgb": "orgb"}}`)),
	}}

	if err := DB.Create(&users).Error; err != nil {
		t.Errorf("Failed to create users %v", err)
	}

	var result UserWithJSON
	if err := DB.First(&result, datatypes.JSONQuery("attributes").HasKey("role")).Error; err != nil {
		t.Fatalf("failed to find user with json key, got error %v", err)
	}
	AssertEqual(t, result.Name, users[1].Name)

	var result2 UserWithJSON
	if err := DB.First(&result2, datatypes.JSONQuery("attributes").HasKey("orgs", "orga")).Error; err != nil {
		t.Fatalf("failed to find user with json key, got error %v", err)
	}
	AssertEqual(t, result2.Name, users[0].Name)
}
