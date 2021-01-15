package datatypes_test

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	. "gorm.io/gorm/utils/tests"
)

var _ driver.Valuer = &datatypes.JSON{}

func TestJSONMap(t *testing.T) {
	if SupportedDriver("sqlite", "mysql", "postgres") {
		type UserWithJSONMap struct {
			gorm.Model
			Name       string
			Attributes datatypes.JSONMap
		}

		DB.Migrator().DropTable(&UserWithJSONMap{})
		if err := DB.Migrator().AutoMigrate(&UserWithJSONMap{}); err != nil {
			t.Errorf("failed to migrate, got error: %v", err)
		}

		// Go's json marshaler removes whitespace & orders keys alphabetically
		// use to compare against marshaled []byte of datatypes.JSON
		user1AttrsStr := `{"age":18,"name":"json-1","orgs":{"orga":"orga"},"tags":["tag1","tag2"]}`
		user1Attrs := map[string]interface{}{
			"age":  18,
			"name": "json-1",
			"orgs": map[string]interface{}{
				"orga": "orga",
			},
			"tags": []interface{}{"tag1", "tag2"},
		}

		user2Attrs := map[string]interface{}{
			"name": "json-2",
			"age":  28,
			"tags": []interface{}{"tag1", "tag3"},
			"role": "admin",
			"orgs": map[string]interface{}{
				"orgb": "orgb",
			},
		}

		users := []UserWithJSONMap{{
			Name:       "json-1",
			Attributes: datatypes.JSONMap(user1Attrs),
		}, {
			Name:       "json-2",
			Attributes: datatypes.JSONMap(user2Attrs),
		}}

		if err := DB.Create(&users).Error; err != nil {
			t.Errorf("Failed to create users %v", err)
		}

		var result UserWithJSONMap
		if err := DB.First(&result, datatypes.JSONQuery("attributes").HasKey("role")).Error; err != nil {
			t.Fatalf("failed to find user with json key, got error %v", err)
		}
		AssertEqual(t, result.Name, users[1].Name)

		var result2 UserWithJSONMap
		if err := DB.First(&result2, datatypes.JSONQuery("attributes").HasKey("orgs", "orga")).Error; err != nil {
			t.Fatalf("failed to find user with json key, got error %v", err)
		}
		AssertEqual(t, result2.Name, users[0].Name)

		//https://tip.golang.org/doc/go1.12#fmt
		AssertEqual(t, fmt.Sprint(result2.Attributes), fmt.Sprint(user1Attrs))

		// attributes should not marshal to base64 encoded []byte
		result2Attrs, err := json.Marshal(result2.Attributes)
		if err != nil {
			t.Fatalf("failed to marshal result2.Attributes, got error %v", err)
		}

		AssertEqual(t, string(result2Attrs), user1AttrsStr)

		// []byte should unmarshal into type datatypes.JSONMap
		var j datatypes.JSONMap
		if err := json.Unmarshal([]byte(user1AttrsStr), &j); err != nil {
			t.Fatalf("failed to unmarshal user1Attrs, got error %v", err)
		}

		AssertEqual(t, fmt.Sprint(j), fmt.Sprint(user1Attrs))

		var result3 UserWithJSONMap
		if err := DB.First(&result3, datatypes.JSONQuery("attributes").Equals("json-1", "name")).Error; err != nil {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}
		AssertEqual(t, result3.Name, users[0].Name)

		var result4 UserWithJSONMap
		if err := DB.First(&result4, datatypes.JSONQuery("attributes").Equals("orgb", "orgs", "orgb")).Error; err != nil {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}
		AssertEqual(t, result4.Name, users[1].Name)

		// FirstOrCreate
		jsonMap := map[string]interface{}{"Attributes": datatypes.JSON(`{"age":19,"name":"json-1","orgs":{"orga":"orga"},"tags":["tag1","tag2"]}`)}
		if err := DB.Where(&UserWithJSONMap{Name: "json-1"}).Assign(jsonMap).FirstOrCreate(&UserWithJSONMap{}).Error; err != nil {
			t.Errorf("failed to run FirstOrCreate")
		}

		var result5 UserWithJSONMap
		if err := DB.First(&result5, datatypes.JSONQuery("attributes").Equals(19, "age")).Error; err != nil {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}
	}
}
