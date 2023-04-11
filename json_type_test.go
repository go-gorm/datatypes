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
		}, {
			Name:       "json-4",
			Attributes: datatypes.NewJSONType(Attribute{Tags: []string{"tag1", "tag2", "tag3"}}),
		}, {
			Name:       "json-5",
			Attributes: datatypes.NewJSONType(Attribute{Tags: []string{"tag1", "tag2", "tag3"}}),
		}}

		if err := DB.Create(&users).Error; err != nil {
			t.Errorf("Failed to create users %v", err)
		}

		var result UserWithJSON
		if err := DB.First(&result, users[1].ID).Error; err != nil {
			t.Fatalf("failed to find user with json key, got error %v", err)
		}
		AssertEqual(t, result.Name, users[1].Name)
		AssertEqual(t, result.Attributes.Data().Age, users[1].Attributes.Data().Age)
		AssertEqual(t, result.Attributes.Data().Admin, users[1].Attributes.Data().Admin)
		AssertEqual(t, len(result.Attributes.Data().Orgs), len(users[1].Attributes.Data().Orgs))

		// List
		var users2 []UserWithJSON
		if err := DB.Model(&UserWithJSON{}).Limit(10).Order("id asc").Find(&users2).Error; err != nil {
			t.Fatalf("failed to select attribute field, got error %v", err)
		}
		AssertEqual(t, users2[0].Attributes.Data().Age, 18)

		// Select Field
		var singleUser UserWithJSON
		if err := DB.Model(&UserWithJSON{}).Select("attributes").Limit(1).Order("id asc").Find(&singleUser).Error; err != nil {
			t.Fatalf("failed to select attribute field, got error %v", err)
		}
		AssertEqual(t, singleUser.Attributes.Data().Age, 18)

		// Pluck
		var attr datatypes.JSONType[Attribute]
		if err := DB.Model(&UserWithJSON{}).Limit(1).Order("id asc").Pluck("attributes", &attr).Error; err != nil {
			t.Fatalf("failed to pluck for field, got error %v", err)
		}
		var attribute = attr.Data()
		AssertEqual(t, attribute.Age, 18)

		// Smart Select Fields
		var row struct {
			Attributes datatypes.JSONType[Attribute]
		}
		if err := DB.Model(&UserWithJSON{}).Limit(1).Order("id asc").Find(&row).Error; err != nil {
			t.Fatalf("failed to select attribute field, got error %v", err)
		}
		AssertEqual(t, row.Attributes.Data().Age, 18)

		// FirstOrCreate
		jsonMap := UserWithJSON{
			Attributes: newJSONType[Attribute]([]byte(`{"age":19,"name":"json-1","orgs":{"orga":"orga"},"tags":["tag1","tag2"]}`)),
		}
		if err := DB.Where(&UserWithJSON{Name: "json-1"}).Assign(jsonMap).FirstOrCreate(&UserWithJSON{}).Error; err != nil {
			t.Errorf("failed to run FirstOrCreate")
		}

		// Update
		jsonMap = UserWithJSON{
			Attributes: datatypes.NewJSONType(
				Attribute{
					Age:  18,
					Sex:  1,
					Orgs: map[string]string{"orga": "orga"},
					Tags: []string{"tag1", "tag2", "tag3"},
				},
			),
		}
		var result3 UserWithJSON
		result3.ID = 1
		if err := DB.Model(&result3).Updates(jsonMap).Error; err != nil {
			t.Errorf("failed to run FirstOrCreate")
		}
	}
}

func TestJSONSlice(t *testing.T) {
	if SupportedDriver("sqlite", "mysql", "postgres") {
		type Tag struct {
			Name  string
			Score float64
		}
		type UserWithJSON2 struct {
			gorm.Model
			Name string
			Tags datatypes.JSONSlice[Tag]
		}
		type UserWithJSON = UserWithJSON2

		DB.Migrator().DropTable(&UserWithJSON{})
		if err := DB.Migrator().AutoMigrate(&UserWithJSON{}); err != nil {
			t.Errorf("failed to migrate, got error: %v", err)
		}

		// Go's json marshaler removes whitespace & orders keys alphabetically
		// use to compare against marshaled []byte of datatypes.JSON
		var tags = []Tag{{Name: "tag1", Score: 0.1}, {Name: "tag2", Score: 0.2}}

		users := []UserWithJSON{{
			Name: "json-1",
			Tags: datatypes.JSONSlice[Tag]{{Name: "tag1", Score: 1.1}, {Name: "tag2", Score: 1.2}},
		}, {
			Name: "json-2",
			Tags: datatypes.NewJSONSlice([]Tag{{Name: "tag3", Score: 0.3}, {Name: "tag4", Score: 0.4}}),
		}, {
			Name: "json-3",
			Tags: datatypes.JSONSlice[Tag](tags),
		}, {
			Name: "json-4",
			Tags: datatypes.NewJSONSlice(tags),
		}}

		if err := DB.Create(&users).Error; err != nil {
			t.Errorf("Failed to create users %v", err)
		}

		var result UserWithJSON
		if err := DB.First(&result, users[0].ID).Error; err != nil {
			t.Fatalf("failed to find user with json key, got error %v", err)
		}
		AssertEqual(t, result.Name, users[0].Name)
		AssertEqual(t, result.Tags[0], users[0].Tags[0])

		// Pluck
		/*
			var pluckTags datatypes.JSONSlice[Tag]
			if err := DB.Model(&UserWithJSON{}).Limit(1).Order("id asc").Pluck("tags", &pluckTags).Error; err != nil {
				t.Fatalf("failed to pluck for field, got error %v", err)
			}
			AssertEqual(t, len(pluckTags), 2)
			AssertEqual(t, pluckTags[0].Name, "tag1")
		*/

		// Smart Select Fields
		var row struct {
			Tags datatypes.JSONSlice[Tag]
		}
		if err := DB.Model(&UserWithJSON{}).Limit(1).Order("id asc").Find(&row).Error; err != nil {
			t.Fatalf("failed to select attribute field, got error %v", err)
		}
		AssertEqual(t, len(row.Tags), 2)
		AssertEqual(t, row.Tags[0].Name, "tag1")

		// FirstOrCreate
		jsonMap := UserWithJSON{
			Tags: datatypes.NewJSONSlice(tags),
		}
		if err := DB.Where(&UserWithJSON{Name: "json-1"}).Assign(jsonMap).FirstOrCreate(&UserWithJSON{}).Error; err != nil {
			t.Errorf("failed to run FirstOrCreate")
		}

		// Update
		jsonMap = UserWithJSON{
			Tags: datatypes.NewJSONSlice(tags),
		}
		var result3 UserWithJSON
		result3.ID = 1
		if err := DB.Model(&result3).Updates(jsonMap).Error; err != nil {
			t.Errorf("failed to run Updates")
		}
	}
}
