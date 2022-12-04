# GORM Data Types

## JSON

sqlite, mysql, postgres supported

```go
import "gorm.io/datatypes"

type UserWithJSON struct {
	gorm.Model
	Name       string
	Attributes datatypes.JSON
}

DB.Create(&User{
	Name:       "json-1",
	Attributes: datatypes.JSON([]byte(`{"name": "jinzhu", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
}

// Check JSON has keys
datatypes.JSONQuery("attributes").HasKey(value, keys...)

db.Find(&user, datatypes.JSONQuery("attributes").HasKey("role"))
db.Find(&user, datatypes.JSONQuery("attributes").HasKey("orgs", "orga"))
// MySQL
// SELECT * FROM `users` WHERE JSON_EXTRACT(`attributes`, '$.role') IS NOT NULL
// SELECT * FROM `users` WHERE JSON_EXTRACT(`attributes`, '$.orgs.orga') IS NOT NULL

// PostgreSQL
// SELECT * FROM "user" WHERE "attributes"::jsonb ? 'role'
// SELECT * FROM "user" WHERE "attributes"::jsonb -> 'orgs' ? 'orga'


// Check JSON extract value from keys equal to value
datatypes.JSONQuery("attributes").Equals(value, keys...)

DB.First(&user, datatypes.JSONQuery("attributes").Equals("jinzhu", "name"))
DB.First(&user, datatypes.JSONQuery("attributes").Equals("orgb", "orgs", "orgb"))
// MySQL
// SELECT * FROM `user` WHERE JSON_EXTRACT(`attributes`, '$.name') = "jinzhu"
// SELECT * FROM `user` WHERE JSON_EXTRACT(`attributes`, '$.orgs.orgb') = "orgb"

// PostgreSQL
// SELECT * FROM "user" WHERE json_extract_path_text("attributes"::json,'name') = 'jinzhu'
// SELECT * FROM "user" WHERE json_extract_path_text("attributes"::json,'orgs','orgb') = 'orgb'
```

NOTE: SQlite need to build with `json1` tag, e.g: `go build --tags json1`, refer https://github.com/mattn/go-sqlite3#usage

## Date

```go
import "gorm.io/datatypes"

type UserWithDate struct {
	gorm.Model
	Name string
	Date datatypes.Date
}

user := UserWithDate{Name: "jinzhu", Date: datatypes.Date(time.Now())}
DB.Create(&user)
// INSERT INTO `user_with_dates` (`name`,`date`) VALUES ("jinzhu","2020-07-17 00:00:00")

DB.First(&result, "name = ? AND date = ?", "jinzhu", datatypes.Date(curTime))
// SELECT * FROM user_with_dates WHERE name = "jinzhu" AND date = "2020-07-17 00:00:00" ORDER BY `user_with_dates`.`id` LIMIT 1
```

## Time

MySQL, PostgreSQL, SQLite, SQLServer are supported.

Time with nanoseconds is supported for some databases which support for time with fractional second scale.

```go
import "gorm.io/datatypes"

type UserWithTime struct {
    gorm.Model
    Name string
    Time datatypes.Time
}

user := UserWithTime{Name: "jinzhu", Time: datatypes.NewTime(1, 2, 3, 0)}
DB.Create(&user)
// INSERT INTO `user_with_times` (`name`,`time`) VALUES ("jinzhu","01:02:03")

DB.First(&result, "name = ? AND time = ?", "jinzhu", datatypes.NewTime(1, 2, 3, 0))
// SELECT * FROM user_with_times WHERE name = "jinzhu" AND time = "01:02:03" ORDER BY `user_with_times`.`id` LIMIT 1
```

NOTE: If the current using database is SQLite, the field column type is defined as `TEXT` type
when GORM AutoMigrate because SQLite doesn't have time type.

## JSON_SET

sqlite, mysql supported

```go
import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserWithJSON struct {
	gorm.Model
	Name       string
	Attributes datatypes.JSON
}

DB.Create(&UserWithJSON{
	Name:       "json-1",
	Attributes: datatypes.JSON([]byte(`{"name": "json-1", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
})

type User struct {
	Name string
	Age  int
}

friend := User{
	Name: "Bob",
	Age:  21,
}

// Set fields of JSON column
datatypes.JSONSet("attributes").Set("age", 20).Set("tags[0]", "tag2").Set("orgs.orga", "orgb")

DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("age", 20).Set("tags[0]", "tag3").Set("orgs.orga", "orgb"))
DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("phones", []string{"10085", "10086"}))
DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("phones", gorm.Expr("CAST(? AS JSON)", `["10085", "10086"]`)))
DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("friend", friend))
// MySQL
// UPDATE `user_with_jsons` SET `attributes` = JSON_SET(`attributes`, '$.tags[0]', 'tag3', '$.orgs.orga', 'orgb', '$.age', 20) WHERE name = 'json-1'
// UPDATE `user_with_jsons` SET `attributes` = JSON_SET(`attributes`, '$.phones', CAST('["10085", "10086"]' AS JSON)) WHERE name = 'json-1'
// UPDATE `user_with_jsons` SET `attributes` = JSON_SET(`attributes`, '$.phones', CAST('["10085", "10086"]' AS JSON)) WHERE name = 'json-1'
// UPDATE `user_with_jsons` SET `attributes` = JSON_SET(`attributes`, '$.friend', CAST('{"Name": "Bob", "Age": 21}' AS JSON)) WHERE name = 'json-1'
```
NOTE: MariaDB does not support CAST(? AS JSON).

## JSONType[T]

sqlite, mysql, postgres supported

```go
import "gorm.io/datatypes"

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

var user = UserWithJSON{
	Name: "hello"
	Attributes: datatypes.JSONType[Attribute]{
		Data: Attribute{
			Age:  18,
			Sex:  1,
			Orgs: map[string]string{"orga": "orga"},
			Tags: []string{"tag1", "tag2", "tag3"},
		},
	},
}

// Create
DB.Create(&user)

// First
var result UserWithJSON
DB.First(&result, user.ID)

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

DB.Model(&user).Updates(jsonMap)
```

NOTE: it's not support json query
