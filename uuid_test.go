package datatypes_test

import (
	"database/sql/driver"
	"testing"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	. "gorm.io/gorm/utils/tests"
)

var _ driver.Valuer = &datatypes.UUID{}

func TestUUID(t *testing.T) {
	if SupportedDriver("sqlite", "mysql", "postgres", "sqlserver") {
		type UserWithUUID struct {
			gorm.Model
			Name     string
			UserUUID datatypes.UUID
		}

		DB.Migrator().DropTable(&UserWithUUID{})
		if err := DB.Migrator().AutoMigrate(&UserWithUUID{}); err != nil {
			t.Errorf("failed to migrate, got error: %v", err)
		}

		user1UUID, uuidGenErr := uuid.NewUUID()
		if uuidGenErr != nil {
			t.Fatalf("failed to generate a new UUID, got error %v", uuidGenErr)
		}
		user2UUID, uuidGenErr := uuid.NewUUID()
		if uuidGenErr != nil {
			t.Fatalf("failed to generate a new UUID, got error %v", uuidGenErr)
		}
		user3UUID, uuidGenErr := uuid.NewUUID()
		if uuidGenErr != nil {
			t.Fatalf("failed to generate a new UUID, got error %v", uuidGenErr)
		}
		user4UUID, uuidGenErr := uuid.NewUUID()
		if uuidGenErr != nil {
			t.Fatalf("failed to generate a new UUID, got error %v", uuidGenErr)
		}

		users := []UserWithUUID{{
			Name:     "uuid-1",
			UserUUID: datatypes.UUID(user1UUID),
		}, {
			Name:     "uuid-2",
			UserUUID: datatypes.UUID(user2UUID),
		}, {
			Name:     "uuid-3",
			UserUUID: datatypes.UUID(user3UUID),
		}, {
			Name:     "uuid-4",
			UserUUID: datatypes.UUID(user4UUID),
		}}

		if err := DB.Create(&users).Error; err != nil {
			t.Errorf("Failed to create users %v", err)
		}

		for _, user := range users {
			result := UserWithUUID{}
			if err := DB.First(
				&result, "name = ? AND user_uuid = ?",
				user.Name,
				user.UserUUID,
			).Error; err != nil {
				t.Fatalf("failed to find user with uuid, got error: %v", err)
			}
			AssertEqual(t, !result.UserUUID.IsEmpty(), true)
			AssertEqual(t, user.UserUUID.Equals(result.UserUUID), true)
			valueUser, err := user.UserUUID.Value()
			if err != nil {
				t.Fatalf("failed to get user value, got error: %v", err)
			}
			valueResult, err := result.UserUUID.Value()
			if err != nil {
				t.Fatalf("failed to get result value, got error: %v", err)
			}
			AssertEqual(t, valueUser, valueResult)
		}

		user1 := users[0]
		AssertEqual(t, user1.UserUUID.IsNil(), false)
		AssertEqual(t, user1.UserUUID.IsEmpty(), false)
		DB.Model(&user1).Updates(
			map[string]interface{}{"user_uuid": uuid.Nil},
		)
		AssertEqual(t, user1.UserUUID.IsNil(), true)
		AssertEqual(t, user1.UserUUID.IsEmpty(), true)

		user2 := users[1]
		AssertEqual(t, user2.UserUUID.IsNil(), false)
		AssertEqual(t, user2.UserUUID.IsEmpty(), false)
		DB.Model(&user2).Updates(
			map[string]interface{}{"user_uuid": nil},
		)
		AssertEqual(t, user2.UserUUID.IsNil(), true)
		AssertEqual(t, user2.UserUUID.IsEmpty(), true)
	}
}

func TestUUIDPtr(t *testing.T) {
	if SupportedDriver("sqlite", "mysql", "postgres", "sqlserver") {
		type UserWithUUIDPtr struct {
			gorm.Model
			Name     string
			UserUUID *datatypes.UUID
		}

		DB.Migrator().DropTable(&UserWithUUIDPtr{})
		if err := DB.Migrator().AutoMigrate(&UserWithUUIDPtr{}); err != nil {
			t.Errorf("failed to migrate, got error: %v", err)
		}

		user1UUID, uuidGenErr := uuid.NewUUID()
		if uuidGenErr != nil {
			t.Fatalf("failed to generate a new UUID, got error %v", uuidGenErr)
		}
		user2UUID, uuidGenErr := uuid.NewUUID()
		if uuidGenErr != nil {
			t.Fatalf("failed to generate a new UUID, got error %v", uuidGenErr)
		}
		user3UUID, uuidGenErr := uuid.NewUUID()
		if uuidGenErr != nil {
			t.Fatalf("failed to generate a new UUID, got error %v", uuidGenErr)
		}
		user4UUID, uuidGenErr := uuid.NewUUID()
		if uuidGenErr != nil {
			t.Fatalf("failed to generate a new UUID, got error %v", uuidGenErr)
		}

		uuid1 := datatypes.UUID(user1UUID)
		uuid2 := datatypes.UUID(user2UUID)
		uuid3 := datatypes.UUID(user3UUID)
		uuid4 := datatypes.UUID(user4UUID)

		users := []UserWithUUIDPtr{{
			Name:     "uuid-1",
			UserUUID: &uuid1,
		}, {
			Name:     "uuid-2",
			UserUUID: &uuid2,
		}, {
			Name:     "uuid-3",
			UserUUID: &uuid3,
		}, {
			Name:     "uuid-4",
			UserUUID: &uuid4,
		}}

		if err := DB.Create(&users).Error; err != nil {
			t.Errorf("Failed to create users %v", err)
		}

		for _, user := range users {
			result := UserWithUUIDPtr{}
			if err := DB.First(
				&result, "name = ? AND user_uuid = ?",
				user.Name,
				*user.UserUUID,
			).Error; err != nil {
				t.Fatalf("failed to find user with uuid, got error: %v", err)
			}
			AssertEqual(t, !result.UserUUID.IsEmpty(), true)
			AssertEqual(t, user.UserUUID, result.UserUUID)
			valueUser, err := user.UserUUID.Value()
			if err != nil {
				t.Fatalf("failed to get user value, got error: %v", err)
			}
			valueResult, err := result.UserUUID.Value()
			if err != nil {
				t.Fatalf("failed to get result value, got error: %v", err)
			}
			AssertEqual(t, valueUser, valueResult)
		}

		user1 := users[0]
		AssertEqual(t, user1.UserUUID.IsNilPtr(), false)
		AssertEqual(t, user1.UserUUID.IsEmptyPtr(), false)
		DB.Model(&user1).Updates(map[string]interface{}{"user_uuid": nil})
		AssertEqual(t, user1.UserUUID.IsNilPtr(), true)
		AssertEqual(t, user1.UserUUID.IsEmptyPtr(), true)
	}
}
