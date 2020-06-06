# GORM Data Types

## JSON

mysql, postgres supported

```go
type UserWithJSON struct {
	gorm.Model
	Name       string
	Attributes datatypes.JSON
}

DB.Create(&User{
	Name:       "json-1",
	Attributes: datatypes.JSON([]byte(`{"name": "jinzhu", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
}

DB.First(&user, datatypes.JSONQuery("json").HasKey("role"))
DB.First(&user, datatypes.JSONQuery("json").HasKey("orgs", "orga"))
```
