package datatypes_test

import (
	"errors"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	. "gorm.io/gorm/utils/tests"
)

func TestNanoID(t *testing.T) {
	if SupportedDriver("postgres") {
		type UserWithNanoID struct {
			ID   datatypes.NanoID `gorm:"nanoid" json:"id"`
			Name string
		}

		DB.Migrator().DropTable(&UserWithNanoID{})
		if err := DB.Migrator().AutoMigrate(&UserWithNanoID{}); err != nil {
			t.Errorf("Failed to migrate, got error: %v", err)
		}

		user := UserWithNanoID{Name: "nkvi.dev"}
		DB.Create(&user)

		result := UserWithNanoID{}
		if err := DB.First(&result, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("Failed to find user with id")
		}

		AssertEqual(t, result.ID, user.ID)
		AssertEqual(t, result.Name, user.Name)

		DB.Delete(&UserWithNanoID{}, "id = ?", result.ID)

		resultAfterDelete := UserWithNanoID{}
		err := DB.First(&resultAfterDelete, "id = ?", user.ID).Error

		AssertEqual(t, errors.Is(err, gorm.ErrRecordNotFound), true)
	}
}
