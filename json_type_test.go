package datatypes_test

import (
	"database/sql/driver"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	. "gorm.io/gorm/utils/tests"
)

var _ driver.Valuer = &datatypes.JSONType[[]int]{}

func newJSONType[T any](b []byte) datatypes.JSONType[T] {
	var t datatypes.JSONType[T]
	_ = t.UnmarshalJSON(b)
	return t
}

func TestJSONType(t *testing.T) {
	if SupportedDriver("sqlite", "mysql", "postgres") {
		type Attribute struct {
			Sex   int
			Age   int
			Orgs  map[string]string
			Tags  []string
			Admin bool
			Role  string
		}
		type UserWithJSON struct {
			gorm.Model
			Name       string
			Attributes datatypes.JSONType[Attribute]
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
			Attributes: newJSONType[Attribute]([]byte(user1Attrs)),
		}, {
			Name:       "json-2",
			Attributes: newJSONType[Attribute]([]byte(`{"name": "json-2", "age": 28, "tags": ["tag1", "tag3"], "role": "admin", "orgs": {"orgb": "orgb"}}`)),
		}, {
			Name:       "json-3",
			Attributes: newJSONType[Attribute]([]byte(`{"tags": ["tag1","tag2","tag3"]`)),
		}}

		if err := DB.Create(&users).Error; err != nil {
			t.Errorf("Failed to create users %v", err)
		}

		var result UserWithJSON
		if err := DB.First(&result, users[1].ID).Error; err != nil {
			t.Fatalf("failed to find user with json key, got error %v", err)
		}
		AssertEqual(t, result.Name, users[1].Name)
		AssertEqual(t, result.Attributes.Data.Age, users[1].Attributes.Data.Age)
		AssertEqual(t, result.Attributes.Data.Admin, users[1].Attributes.Data.Admin)
		AssertEqual(t, len(result.Attributes.Data.Orgs), len(users[1].Attributes.Data.Orgs))

		// FirstOrCreate
		jsonMap := UserWithJSON{
			Attributes: newJSONType[Attribute]([]byte(`{"age":19,"name":"json-1","orgs":{"orga":"orga"},"tags":["tag1","tag2"]}`)),
		}
		if err := DB.Where(&UserWithJSON{Name: "json-1"}).Assign(jsonMap).FirstOrCreate(&UserWithJSON{}).Error; err != nil {
			t.Errorf("failed to run FirstOrCreate")
		}

		// Update
		jsonMap = UserWithJSON{
			Attributes: datatypes.JSONType[Attribute]{
				Data: Attribute{
					Age:  18,
					Sex:  1,
					Orgs: map[string]string{"orga": "orga"},
					Tags: []string{"tag1", "tag2", "tag3"},
				},
			},
		}
		var result3 UserWithJSON
		result3.ID = 1
		if err := DB.Model(&result3).Updates(jsonMap).Error; err != nil {
			t.Errorf("failed to run FirstOrCreate")
		}
	}
}
