package datatypes_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	. "gorm.io/gorm/utils/tests"
)

var _ driver.Valuer = &datatypes.JSON{}

func TestJSON(t *testing.T) {
	if SupportedDriver("sqlite", "mysql", "postgres") {
		type UserWithJSON struct {
			gorm.Model
			Name       string
			Attributes datatypes.JSON
		}

		DB.Migrator().DropTable(&UserWithJSON{})
		if err := DB.Migrator().AutoMigrate(&UserWithJSON{}); err != nil {
			t.Errorf("failed to migrate, got error: %v", err)
		}

		// Go's json marshaler removes whitespace & orders keys alphabetically
		// use to compare against marshaled []byte of datatypes.JSON
		user1Attrs := `{"age":18,"name":"json-1","orgs":{"orga":"orga"},"tags":["tag1","tag2"],"admin":true}`

		users := []UserWithJSON{{
			Name:       "json-1",
			Attributes: datatypes.JSON([]byte(user1Attrs)),
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

		// attributes should not marshal to base64 encoded []byte
		result2Attrs, err := json.Marshal(&result2.Attributes)
		if err != nil {
			t.Fatalf("failed to marshal result2.Attributes, got error %v", err)
		}
		AssertEqual(t, string(result2Attrs), user1Attrs)

		// []byte should unmarshal into type datatypes.JSON
		var j datatypes.JSON
		if err := json.Unmarshal([]byte(user1Attrs), &j); err != nil {
			t.Fatalf("failed to unmarshal user1Attrs, got error %v", err)
		}

		AssertEqual(t, string(j), user1Attrs)

		var result3 UserWithJSON
		if err := DB.First(&result3, datatypes.JSONQuery("attributes").Equals("json-1", "name")).Error; err != nil {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}
		AssertEqual(t, result3.Name, users[0].Name)

		var result4 UserWithJSON
		if err := DB.First(&result4, datatypes.JSONQuery("attributes").Equals("orgb", "orgs", "orgb")).Error; err != nil {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}
		AssertEqual(t, result4.Name, users[1].Name)

		var results5 []UserWithJSON
		if err := DB.Where(datatypes.JSONQuery("attributes").HasKey("age")).Where(datatypes.JSONQuery("attributes").Equals(true, "admin")).Find(&results5).Error; err != nil || len(results5) != 1 {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}
		AssertEqual(t, results5[0].Name, users[0].Name)

		// FirstOrCreate
		jsonMap := map[string]interface{}{"Attributes": datatypes.JSON(`{"age":19,"name":"json-1","orgs":{"orga":"orga"},"tags":["tag1","tag2"]}`)}
		if err := DB.Where(&UserWithJSON{Name: "json-1"}).Assign(jsonMap).FirstOrCreate(&UserWithJSON{}).Error; err != nil {
			t.Errorf("failed to run FirstOrCreate")
		}

		var result6 UserWithJSON
		if err := DB.First(&result6, datatypes.JSONQuery("attributes").Equals(19, "age")).Error; err != nil {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}

		// Update
		jsonMap = map[string]interface{}{"Attributes": datatypes.JSON(`{"age":29,"name":"json-1","orgs":{"orga":"orga"},"tags":["tag1","tag2"]}`)}
		if err := DB.Model(&result3).Updates(jsonMap).Error; err != nil {
			t.Errorf("failed to run FirstOrCreate")
		}

		var result7 UserWithJSON
		if err := DB.First(&result7, datatypes.JSONQuery("attributes").Equals(29, "age")).Error; err != nil {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}
	}
}
