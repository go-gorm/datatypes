package datatypes_test

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
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

		var result10 UserWithJSON
		if err := DB.First(&result10, datatypes.JSONQuery("attributes").Likes("%dmi%", "role")).Error; err != nil {
			t.Fatalf("failed to find user with json value, got error %v", err)
		}
		AssertEqual(t, result10.Name, users[1].Name)

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
		if err := DB.Where("id = ?", cmp1.ID).First(&retSingle1).Error; err != nil {
			t.Errorf("Failed to find param %v", err)
		}

		var retSingle2 Param
		if err := DB.Where("id = ?", cmp2.ID).First(&retSingle2).Error; err != nil {
			t.Errorf("Failed to find param %v", err)
		}

		AssertEqual(t, retSingle1, cmp1)
		AssertEqual(t, retSingle2, cmp2)

		var retMultiple []Param
		if err := DB.Find(&retMultiple).Error; err != nil {
			t.Errorf("Failed to find param %v", err)
		}

		AssertEqual(t, len(retMultiple), 2)
		AssertEqual(t, retMultiple[0], cmp1)
		AssertEqual(t, retMultiple[1], cmp2)
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

func TestJSONArrayAppend(t *testing.T) {
	if !SupportedDriver("sqlite", "mysql", "postgres") {
		t.Skip()
	}

	type UserWithJSON struct {
		gorm.Model
		Name       string
		Attributes datatypes.JSON
	}

	type Friend struct {
		Name string
		Age  int
	}

	// assertAppend creates a row with attributes, appends values to it and
	// compares the stored JSON with expect. An empty expect means the column
	// must still be NULL.
	assertAppend := func(t *testing.T, model interface{}, name string, attributes datatypes.JSON, values []interface{}, expect string) {
		t.Helper()

		var count int64
		if err := DB.Model(model).Where("name = ? AND attributes IS NULL", name).Count(&count).Error; err != nil {
			t.Fatalf("failed to count users, got error %v", err)
		}
		if (attributes == nil) != (count == 1) {
			t.Fatalf("expect attributes NULL to be %v before append", attributes == nil)
		}

		if err := DB.Model(model).Where("name = ?", name).UpdateColumn("attributes", datatypes.JSONArrayAppend("attributes", values...)).Error; err != nil {
			t.Fatalf("failed to append to json array, got error %v", err)
		}

		if expect == "" {
			if err := DB.Model(model).Where("name = ? AND attributes IS NULL", name).Count(&count).Error; err != nil {
				t.Fatalf("failed to count users, got error %v", err)
			}
			if count != 1 {
				t.Errorf("expect attributes to stay NULL")
			}
			return
		}

		var raw string
		if err := DB.Model(model).Select("attributes").Where("name = ?", name).Scan(&raw).Error; err != nil {
			t.Fatalf("failed to find user, got error %v", err)
		}
		var actual, expected interface{}
		if err := json.Unmarshal([]byte(raw), &actual); err != nil {
			t.Fatalf("failed to unmarshal attributes %q, got error %v", raw, err)
		}
		if err := json.Unmarshal([]byte(expect), &expected); err != nil {
			t.Fatalf("failed to unmarshal expect %q, got error %v", expect, err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expect attributes %s, got %s", expect, raw)
		}
	}

	DB.Migrator().DropTable(&UserWithJSON{})
	if err := DB.Migrator().AutoMigrate(&UserWithJSON{}); err != nil {
		t.Fatalf("failed to migrate, got error: %v", err)
	}

	var jsonExpr interface{}
	switch DB.Dialector.Name() {
	case "mysql":
		jsonExpr = gorm.Expr("JSON_EXTRACT(?,'$')", `["b","c"]`)
	case "sqlite":
		jsonExpr = gorm.Expr("JSON(?)", `["b","c"]`)
	case "postgres":
		jsonExpr = gorm.Expr("?::jsonb", `["b","c"]`)
	}

	tests := []struct {
		name       string
		attributes datatypes.JSON // nil is stored as NULL
		values     []interface{}
		expect     string
	}{
		{"null column", nil, []interface{}{"tag1"}, `["tag1"]`},
		{"empty array", datatypes.JSON(`[]`), []interface{}{"tag1"}, `["tag1"]`},
		{"strings", datatypes.JSON(`["tag1"]`), []interface{}{"tag2", "标签", `a"b\c`}, `["tag1","tag2","标签","a\"b\\c"]`},
		{"json-like strings", datatypes.JSON(`[]`), []interface{}{`[1,2]`, "true", "1", "null"}, `["[1,2]","true","1","null"]`},
		{"numbers", datatypes.JSON(`[1]`), []interface{}{2, 3.5, -4}, `[1,2,3.5,-4]`},
		{"bools", datatypes.JSON(`[]`), []interface{}{true, false}, `[true,false]`},
		{"nil", datatypes.JSON(`[]`), []interface{}{nil}, `[null]`},
		{"map", datatypes.JSON(`[{"k":"v1"}]`), []interface{}{map[string]interface{}{"k": "v2", "n": 1}}, `[{"k":"v1"},{"k":"v2","n":1}]`},
		{"struct", datatypes.JSON(`[]`), []interface{}{Friend{Name: "Bob", Age: 21}}, `[{"Name":"Bob","Age":21}]`},
		{"slice", datatypes.JSON(`[1]`), []interface{}{[]int{2, 3}}, `[1,[2,3]]`},
		{"empty slice", datatypes.JSON(`[1]`), []interface{}{[]string{}}, `[1,[]]`},
		{"nested slice", datatypes.JSON(`[]`), []interface{}{[][]string{{"a"}, {"b", "c"}}}, `[[["a"],["b","c"]]]`},
		{"json", datatypes.JSON(`[]`), []interface{}{datatypes.JSON(`{"tags":["a"]}`)}, `[{"tags":["a"]}]`},
		{"multiple values", nil, []interface{}{"a", 1, true, []string{"b", "c"}, map[string]string{"k": "v"}}, `["a",1,true,["b","c"],{"k":"v"}]`},
		{"expr", datatypes.JSON(`["a"]`), []interface{}{jsonExpr}, `["a",["b","c"]]`},
		{"no values", datatypes.JSON(`["a"]`), nil, `["a"]`},
		{"no values on null column", nil, nil, ""},
	}

	for idx, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name := fmt.Sprintf("json-append-%d", idx)
			if err := DB.Create(&UserWithJSON{Name: name, Attributes: test.attributes}).Error; err != nil {
				t.Fatalf("failed to create user, got error %v", err)
			}
			assertAppend(t, &UserWithJSON{}, name, test.attributes, test.values, test.expect)
		})
	}

	t.Run("invalid value", func(t *testing.T) {
		if err := DB.Create(&UserWithJSON{Name: "json-append-invalid", Attributes: datatypes.JSON(`["a"]`)}).Error; err != nil {
			t.Fatalf("failed to create user, got error %v", err)
		}
		if err := DB.Model(&UserWithJSON{}).Where("name = ?", "json-append-invalid").UpdateColumn("attributes", datatypes.JSONArrayAppend("attributes", make(chan int))).Error; err == nil {
			t.Fatalf("expect error when value can not be marshaled")
		}
		assertAppend(t, &UserWithJSON{}, "json-append-invalid", datatypes.JSON(`["a"]`), nil, `["a"]`)
	})

	if DB.Dialector.Name() == "postgres" {
		t.Run("postgres json column", func(t *testing.T) {
			type UserWithPlainJSON struct {
				gorm.Model
				Name       string
				Attributes datatypes.JSON `gorm:"type:json"`
			}

			DB.Migrator().DropTable(&UserWithPlainJSON{})
			if err := DB.Migrator().AutoMigrate(&UserWithPlainJSON{}); err != nil {
				t.Fatalf("failed to migrate, got error: %v", err)
			}

			users := []UserWithPlainJSON{{Name: "json-1"}, {Name: "json-2", Attributes: datatypes.JSON(`[1]`)}}
			if err := DB.Create(&users).Error; err != nil {
				t.Fatalf("failed to create users, got error %v", err)
			}
			assertAppend(t, &UserWithPlainJSON{}, "json-1", nil, []interface{}{"a", []int{2, 3}}, `["a",[2,3]]`)
			assertAppend(t, &UserWithPlainJSON{}, "json-2", datatypes.JSON(`[1]`), []interface{}{"a", []int{2, 3}}, `[1,"a",[2,3]]`)
		})
	}
}

func TestJSONArrayQuery(t *testing.T) {
	if SupportedDriver("sqlite", "mysql") {
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
		cmp3 := Param{
			DisplayName: "JSONArray-3",
			Config:      datatypes.JSON("{\"test\": [\"a\", \"b\"]}"),
		}
		cmp4 := Param{
			DisplayName: "JSONArray-4",
			Config:      datatypes.JSON("{\"test\": \"c\"}"),
		}

		if err := DB.Create(&cmp1).Error; err != nil {
			t.Errorf("Failed to create param %v", err)
		}
		if err := DB.Create(&cmp2).Error; err != nil {
			t.Errorf("Failed to create param %v", err)
		}
		if err := DB.Create(&cmp3).Error; err != nil {
			t.Errorf("Failed to create param %v", err)
		}
		if err := DB.Create(&cmp4).Error; err != nil {
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

		if err := DB.Where(datatypes.JSONArrayQuery("config").Contains("a", "test")).Find(&retMultiple).Error; err != nil {
			t.Fatalf("failed to find params with json value and keys, got error %v", err)
		}
		AssertEqual(t, len(retMultiple), 1)

		if err := DB.Where(datatypes.JSONArrayQuery("config").In([]string{"c", "a"})).Find(&retMultiple).Error; err != nil {
			t.Fatalf("failed to find params with json value, got error %v", err)
		}
		AssertEqual(t, len(retMultiple), 1)

		if err := DB.Where(datatypes.JSONArrayQuery("config").In([]string{"c", "d"}, "test")).Find(&retMultiple).Error; err != nil {
			t.Fatalf("failed to find params with json value and keys, got error %v", err)
		}
		AssertEqual(t, len(retMultiple), 1)
	}
}

// TestJSONValueTypes tests the return types from JSON Value() methods
func TestJSONValueTypes(t *testing.T) {
	// Test JSON type
	jsonData := datatypes.JSON(`{"test": "value"}`)
	value, err := jsonData.Value()
	if err != nil {
		t.Errorf("JSON.Value() error: %v", err)
	}
	if _, ok := value.(string); !ok {
		t.Errorf("JSON.Value() should return string, got %T", value)
	}

	// Test empty JSON
	emptyJSON := datatypes.JSON("")
	emptyValue, err := emptyJSON.Value()
	if err != nil {
		t.Errorf("Empty JSON.Value() error: %v", err)
	}
	if emptyValue != nil {
		t.Errorf("Empty JSON.Value() should return nil, got %v", emptyValue)
	}
}
