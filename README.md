# GORM Data Types

## JSON

mysql, postgres supported

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

// PostgreSQL
// SELECT * FROM "user" WHERE json_extract_path_text("attributes"::json,'name') = 'jinzhu'
```
