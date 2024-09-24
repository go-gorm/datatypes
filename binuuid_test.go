package datatypes_test

import (
	"database/sql/driver"
	"testing"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	. "gorm.io/gorm/utils/tests"
)

var _ driver.Valuer = &datatypes.BinUUID{}

func TestBinUUID(t *testing.T) {
	if SupportedDriver("sqlite", "mysql", "postgres", "sqlserver") {
		type UserWithBinUUID struct {
			gorm.Model
			Name     string
			UserUUID datatypes.BinUUID
		}

		DB.Migrator().DropTable(&UserWithBinUUID{})
		if err := DB.Migrator().AutoMigrate(&UserWithBinUUID{}); err != nil {
			t.Errorf("failed to migrate, got error: %v", err)
		}

		users := []UserWithBinUUID{{
			Name:     "uuid-1",
			UserUUID: datatypes.NewBinUUIDv1(),
		}, {
			Name:     "uuid-2",
			UserUUID: datatypes.NewBinUUIDv1(),
		}, {
			Name:     "uuid-3",
			UserUUID: datatypes.NewBinUUIDv4(),
		}, {
			Name:     "uuid-4",
			UserUUID: datatypes.NewBinUUIDv4(),
		}}

		if err := DB.Create(&users).Error; err != nil {
			t.Errorf("Failed to create users %v", err)
		}

		for _, user := range users {
			result := UserWithBinUUID{}
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
			AssertEqual(t, user.UserUUID.LengthBytes(), 16)
			AssertEqual(t, user.UserUUID.Length(), 36)
		}

		var tx *gorm.DB
		user1 := users[0]
		AssertEqual(t, user1.UserUUID.IsNil(), false)
		AssertEqual(t, user1.UserUUID.IsEmpty(), false)
		tx = DB.Model(&user1).Updates(
			map[string]interface{}{
				"user_uuid": datatypes.BinUUIDFromString(uuid.Nil.String()),
			},
		)
		AssertEqual(t, tx.Error, nil)
		AssertEqual(t, user1.UserUUID.IsNil(), true)
		AssertEqual(t, user1.UserUUID.IsEmpty(), true)
		user1NewUUID := datatypes.NewBinUUIDv4()
		tx = DB.Model(&user1).Updates(
			map[string]interface{}{
				"user_uuid": user1NewUUID,
			},
		)
		AssertEqual(t, tx.Error, nil)
		AssertEqual(t, user1.UserUUID, user1NewUUID)

		user2 := users[1]
		AssertEqual(t, user2.UserUUID.IsNil(), false)
		AssertEqual(t, user2.UserUUID.IsEmpty(), false)
		tx = DB.Model(&user2).Updates(
			map[string]interface{}{"user_uuid": datatypes.NewNilBinUUID()},
		)
		AssertEqual(t, tx.Error, nil)
		AssertEqual(t, user2.UserUUID.IsNil(), true)
		AssertEqual(t, user2.UserUUID.IsEmpty(), true)
		user2NewUUID := datatypes.NewBinUUIDv4()
		tx = DB.Model(&user2).Updates(
			map[string]interface{}{
				"user_uuid": user2NewUUID,
			},
		)
		AssertEqual(t, tx.Error, nil)
		AssertEqual(t, user2.UserUUID, user2NewUUID)
	}
}

func TestBinUUIDPtr(t *testing.T) {
	if SupportedDriver("sqlite", "mysql", "postgres", "sqlserver") {
		type UserWithBinUUIDPtr struct {
			gorm.Model
			Name     string
			UserUUID *datatypes.BinUUID
		}

		DB.Migrator().DropTable(&UserWithBinUUIDPtr{})
		if err := DB.Migrator().AutoMigrate(&UserWithBinUUIDPtr{}); err != nil {
			t.Errorf("failed to migrate, got error: %v", err)
		}

		uuid1 := datatypes.NewBinUUIDv1()
		uuid2 := datatypes.NewBinUUIDv1()
		uuid3 := datatypes.NewBinUUIDv4()
		uuid4 := datatypes.NewBinUUIDv4()

		users := []UserWithBinUUIDPtr{{
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
			result := UserWithBinUUIDPtr{}
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
			AssertEqual(t, user.UserUUID.LengthBytes(), 16)
			AssertEqual(t, user.UserUUID.Length(), 36)
		}

		user1 := users[0]
		AssertEqual(t, user1.UserUUID.IsNilPtr(), false)
		AssertEqual(t, user1.UserUUID.IsEmptyPtr(), false)
		tx := DB.Model(&user1).Updates(map[string]interface{}{
			"user_uuid": datatypes.NewNilBinUUID(),
		})
		AssertEqual(t, tx.Error, nil)
		AssertEqual(t, user1.UserUUID.IsNil(), true)
		AssertEqual(t, user1.UserUUID.IsEmptyPtr(), true)
	}
}
