package datatypes_test

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
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
		}, {
			Name:       "json-3",
			Attributes: datatypes.JSON([]byte(`["tag1","tag2","tag3"]`)),
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

		var result8 UserWithJSON
		if err := DB.Where(result7, "Attributes").First(&result8).Error; err != nil {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}

		var results9 []UserWithJSON
		if err := DB.Where("? = ?", datatypes.JSONQuery("attributes").Extract("name"), "json-2").Find(&results9).Error; err != nil || len(results9) != 1 {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}
		AssertEqual(t, results9[0].Name, users[1].Name)

		// not support for sqlite
		// JSONOverlaps
		//var result9 UserWithJSON
		//if err := DB.First(&result9, datatypes.JSONOverlaps("attributes", `["tag1","tag2"]`)).Error; err != nil {
		//	t.Fatalf("failed to find user with json value, got error %v", err)
		//}
	}
}

func TestJSONSliceScan(t *testing.T) {
	if SupportedDriver("sqlite", "mysql", "postgres") {
		type Param struct {
			ID          int
			DisplayName string
			Config      datatypes.JSON
		}

		DB.Migrator().DropTable(&Param{})
		if err := DB.Migrator().AutoMigrate(&Param{}); err != nil {
			t.Errorf("failed to migrate, got error: %v", err)
		}

		cmp1 := Param{
			DisplayName: "TestJSONSliceScan-1",
			Config:      datatypes.JSON("{\"param1\": 1234, \"param2\": \"test\"}"),
		}

		cmp2 := Param{
			DisplayName: "TestJSONSliceScan-2",
			Config:      datatypes.JSON("{\"param1\": 456, \"param2\": \"test2\"}"),
		}

		if err := DB.Create(&cmp1).Error; err != nil {
			t.Errorf("Failed to create param %v", err)
		}
		if err := DB.Create(&cmp2).Error; err != nil {
			t.Errorf("Failed to create param %v", err)
		}

		var retSingle1 Param
		if err := DB.Where("id = ?", cmp2.ID).First(&retSingle1).Error; err != nil {
			t.Errorf("Failed to find param %v", err)
		}

		var retSingle2 Param
		if err := DB.Where("id = ?", cmp2.ID).First(&retSingle2).Error; err != nil {
			t.Errorf("Failed to find param %v", err)
		}

		AssertEqual(t, retSingle1, cmp2)
		AssertEqual(t, retSingle2, cmp2)

		var retMultiple []Param
		if err := DB.Find(&retMultiple).Error; err != nil {
			t.Errorf("Failed to find param %v", err)
		}

		AssertEqual(t, retSingle1, cmp2)
		AssertEqual(t, retSingle2, cmp2)
	}
}

func TestPostgresJSONSet(t *testing.T) {
	if !SupportedDriver("postgres") {
		t.Skip()
	}

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
		Attributes: datatypes.JSON(`{"name": "json-1", "age": 18, "orgs": {"orga": "orga"}, "tags": ["tag1", "tag2"], "admin": true}`),
	}, {
		Name:       "json-2",
		Attributes: datatypes.JSON(`{"name": "json-2", "age": 28, "tags": ["tag1", "tag3"], "role": "admin", "orgs": {"orgb": "orgb"}}`),
	}, {
		Name:       "json-3",
		Attributes: datatypes.JSON(`{"name": "json-3"}`),
	}, {
		Name:       "json-4",
		Attributes: datatypes.JSON(`{"name": "json-4"}`),
	}}

	if err := DB.Create(&users).Error; err != nil {
		t.Errorf("Failed to create users %v", err)
	}

	tests := []struct {
		name       string
		userName   string
		path2value map[string]interface{}
		expect     map[string]interface{}
	}{
		{
			name:     "update int and string",
			userName: "json-1",
			path2value: map[string]interface{}{
				"{age}":  20,
				"{role}": "tester",
			},
			expect: map[string]interface{}{
				"age":  20,
				"role": "tester",
			},
		}, {
			name:     "update array child",
			userName: "json-2",
			path2value: map[string]interface{}{
				"{tags, 0}": "tag2",
			},
			expect: map[string]interface{}{
				"tags": []string{"tag2", "tag3"},
			},
		}, {
			name:     "update array",
			userName: "json-2",
			path2value: map[string]interface{}{
				"{phones}": []string{"10086", "10085"},
			},
			expect: map[string]interface{}{
				"phones": []string{"10086", "10085"},
			},
		}, {
			name:     "update by expr",
			userName: "json-4",
			path2value: map[string]interface{}{
				"{extra}": gorm.Expr("?::jsonb", `["a", "b"]`),
			},
			expect: map[string]interface{}{
				"extra": []string{"a", "b"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonSet := datatypes.JSONSet("attributes")
			for path, value := range test.path2value {
				jsonSet = jsonSet.Set(path, value)
			}

			if err := DB.Model(&UserWithJSON{}).Where("name = ?", test.userName).UpdateColumn("attributes", jsonSet).Error; err != nil {
				t.Fatalf("failed to update user with json key, got error %v", err)
			}

			var result UserWithJSON
			if err := DB.First(&result, "name = ?", test.userName).Error; err != nil {
				t.Fatalf("failed to find user with json key, got error %v", err)
			}
			actual := make(map[string]interface{})
			if err := json.Unmarshal(result.Attributes, &actual); err != nil {
				t.Fatalf("failed to unmarshal attributes, got err %v", err)
			}

			for key, value := range test.expect {
				AssertEqual(t, value, test.expect[key])
			}
		})
	}

	if err := DB.Model(&UserWithJSON{}).Where("name = ?", "json-3").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("{friend}", users[0])).Error; err != nil {
		t.Fatalf("failed to update user with json key, got error %v", err)
	}
	var result UserWithJSON
	if err := DB.First(&result, "name = ?", "json-3").Error; err != nil {
		t.Fatalf("failed to find user with json key, got error %v", err)
	}
	actual := make(map[string]json.RawMessage)
	if err := json.Unmarshal(result.Attributes, &actual); err != nil {
		t.Fatalf("failed to unmarshal attributes, got err %v", err)
	}
	var friend UserWithJSON
	if err := json.Unmarshal(actual["friend"], &friend); err != nil {
		t.Fatalf("failed to unmarshal attributes, got err %v", err)
	}
	AssertEqual(t, friend.ID, users[0].ID)
	AssertEqual(t, friend.Name, users[0].Name)
}

func TestJSONSet(t *testing.T) {
	if SupportedDriver("sqlite", "mysql") {
		type UserWithJSON struct {
			gorm.Model
			Name       string
			Attributes datatypes.JSON
		}

		DB.Migrator().DropTable(&UserWithJSON{})
		if err := DB.Migrator().AutoMigrate(&UserWithJSON{}); err != nil {
			t.Errorf("failed to migrate, got error: %v", err)
		}

		var isMariaDB bool
		if DB.Dialector.Name() == "mysql" {
			if v, ok := DB.Dialector.(*mysql.Dialector); ok {
				isMariaDB = strings.Contains(v.ServerVersion, "MariaDB")
			}
		}
		users := []UserWithJSON{{
			Name:       "json-1",
			Attributes: datatypes.JSON([]byte(`{"name": "json-1", "age": 18, "orgs": {"orga": "orga"}, "tags": ["tag1", "tag2"], "admin": true}`)),
		}, {
			Name:       "json-2",
			Attributes: datatypes.JSON([]byte(`{"name": "json-2", "age": 28, "tags": ["tag1", "tag3"], "role": "admin", "orgs": {"orgb": "orgb"}}`)),
		}, {
			Name:       "json-3",
			Attributes: datatypes.JSON([]byte(`{"name": "json-3"}`)),
		}, {
			Name:       "json-4",
			Attributes: datatypes.JSON([]byte(`{"name": "json-4"}`)),
		}}

		if err := DB.Create(&users).Error; err != nil {
			t.Errorf("Failed to create users %v", err)
		}

		tmp := make(map[string]interface{})

		// update int, string
		if err := DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("age", 20).Set("role", "tester")).Error; err != nil {
			t.Fatalf("failed to update user with json key, got error %v", err)
		}
		var result UserWithJSON
		if err := DB.First(&result, "name = ?", "json-1").Error; err != nil {
			t.Fatalf("failed to find user with json key, got error %v", err)
		}
		_ = json.Unmarshal(result.Attributes, &tmp)
		AssertEqual(t, tmp["age"], 20)
		AssertEqual(t, tmp["role"], "tester")

		if err := DB.Model(&UserWithJSON{}).Where("name = ?", "json-2").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("tags[0]", "tag2")).Error; err != nil {
			t.Fatalf("failed to update user with json key, got error %v", err)
		}
		var result2 UserWithJSON
		if err := DB.First(&result2, "name = ?", "json-2").Error; err != nil {
			t.Fatalf("failed to find user with json key, got error %v", err)
		}
		_ = json.Unmarshal(result2.Attributes, &tmp)
		AssertEqual(t, tmp["tags"], []string{"tag2", "tag3"})

		// MariaDB does not support CAST(? AS JSON),
		if isMariaDB {
			if err := DB.Model(&UserWithJSON{}).Where("name = ?", "json-2").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("phones", []string{"10086", "10085"})).Error; err != nil {
				t.Fatalf("failed to update user with json key, got error %v", err)
			}
			var result3 UserWithJSON
			if err := DB.First(&result3, "name = ?", "json-2").Error; err != nil {
				t.Fatalf("failed to find user with json key, got error %v", err)
			}
			_ = json.Unmarshal(result3.Attributes, &tmp)
			AssertEqual(t, tmp["phones"], `["10086","10085"]`)

			if err := DB.Model(&UserWithJSON{}).Where("name = ?", "json-3").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("friend", result)).Error; err != nil {
				t.Fatalf("failed to update user with json key, got error %v", err)
			}
			var result4 UserWithJSON
			if err := DB.First(&result4, "name = ?", "json-3").Error; err != nil {
				t.Fatalf("failed to find user with json key, got error %v", err)
			}
			m := make(map[string]interface{})

			_ = json.Unmarshal(result4.Attributes, &m)
			var tmpResult UserWithJSON
			_ = json.Unmarshal([]byte(m["friend"].(string)), &tmpResult)
			AssertEqual(t, tmpResult.Name, result.Name)

		} else {
			if err := DB.Model(&UserWithJSON{}).Where("name = ?", "json-2").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("phones", []string{"10086", "10085"})).Error; err != nil {
				t.Fatalf("failed to update user with json key, got error %v", err)
			}
			var result3 UserWithJSON
			if err := DB.First(&result3, "name = ?", "json-2").Error; err != nil {
				t.Fatalf("failed to find user with json key, got error %v", err)
			}
			_ = json.Unmarshal(result3.Attributes, &tmp)
			AssertEqual(t, tmp["phones"], []string{"10086", "10085"})

			if err := DB.Model(&UserWithJSON{}).Where("name = ?", "json-3").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("friend", result)).Error; err != nil {
				t.Fatalf("failed to update user with json key, got error %v", err)
			}
			var result4 UserWithJSON
			if err := DB.First(&result4, "name = ?", "json-3").Error; err != nil {
				t.Fatalf("failed to find user with json key, got error %v", err)
			}
			m := make(map[string]UserWithJSON)
			_ = json.Unmarshal(result4.Attributes, &m)
			AssertEqual(t, m["friend"], result)

			if DB.Dialector.Name() == "mysql" {
				if err := DB.Model(&UserWithJSON{}).Where("name = ?", "json-4").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("extra", gorm.Expr("CAST(? AS JSON)", `["a", "b"]`))).Error; err != nil {
					t.Fatalf("failed to update user with json key, got error %v", err)
				}
			} else if DB.Dialector.Name() == "sqlite" {
				if err := DB.Model(&UserWithJSON{}).Where("name = ?", "json-4").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("extra", gorm.Expr("JSON(?)", `["a", "b"]`))).Error; err != nil {
					t.Fatalf("failed to update user with json key, got error %v", err)
				}
			}
			var result5 UserWithJSON
			if err := DB.First(&result5, "name = ?", "json-4").Error; err != nil {
				t.Fatalf("failed to find user with json key, got error %v", err)
			}
			_ = json.Unmarshal(result5.Attributes, &tmp)
			AssertEqual(t, tmp["extra"], []string{"a", "b"})
		}
	}
}

func TestJSONArrayQuery(t *testing.T) {
	if SupportedDriver("mysql") {
		type Param struct {
			ID          int
			DisplayName string
			Config      datatypes.JSON
		}

		DB.Migrator().DropTable(&Param{})
		if err := DB.Migrator().AutoMigrate(&Param{}); err != nil {
			t.Errorf("failed to migrate, got error: %v", err)
		}

		cmp1 := Param{
			DisplayName: "JSONArray-1",
			Config:      datatypes.JSON("[\"a\", \"b\"]"),
		}

		cmp2 := Param{
			DisplayName: "JSONArray-2",
			Config:      datatypes.JSON("[\"c\", \"a\"]"),
		}

		if err := DB.Create(&cmp1).Error; err != nil {
			t.Errorf("Failed to create param %v", err)
		}
		if err := DB.Create(&cmp2).Error; err != nil {
			t.Errorf("Failed to create param %v", err)
		}

		var retSingle1 Param
		if err := DB.Where("id = ?", cmp2.ID).First(&retSingle1).Error; err != nil {
			t.Errorf("Failed to find param %v", err)
		}

		var retSingle2 Param
		if err := DB.Where("id = ?", cmp2.ID).First(&retSingle2).Error; err != nil {
			t.Errorf("Failed to find param %v", err)
		}

		AssertEqual(t, retSingle1, cmp2)
		AssertEqual(t, retSingle2, cmp2)

		var retMultiple []Param

		if err := DB.Where(datatypes.JSONArrayQuery("config").Contains("c")).Find(&retMultiple).Error; err != nil {
			t.Fatalf("failed to find params with json value, got error %v", err)
		}
		AssertEqual(t, len(retMultiple), 1)

	}
}
