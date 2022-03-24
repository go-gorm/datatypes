package datatypes_test

import (
	"net/url"
	"testing"

	"gorm.io/datatypes"
	. "gorm.io/gorm/utils/tests"
)

func TestURL(t *testing.T) {
	type StructWithURL struct {
		ID       uint
		FileName string
		Storage  datatypes.URL
	}
	_ = DB.Migrator().DropTable(&StructWithURL{})
	if err := DB.Migrator().AutoMigrate(&StructWithURL{}); err != nil {
		t.Fatalf("failed to migrate, got error: %v", err)
	}
	f1 := StructWithURL{
		FileName: "FLocal1",
		Storage: datatypes.URL{
			Scheme: "file",
			Path:   "/tmp/f1",
		},
	}
	us := "sftp://user:pwd@127.0.0.1/f2?query=1#frag"
	u2, _ := url.Parse(us)
	f2 := StructWithURL{
		FileName: "FRemote2",
		Storage:  datatypes.URL(*u2),
	}

	uf1 := url.URL(f1.Storage)
	uf2 := url.URL(f2.Storage)
	DB.Create(&f1)
	DB.Create(&f2)

	result := StructWithURL{}
	if err := DB.First(&result, "file_name = ? AND storage = ?", "FLocal1",
		datatypes.URL{
			Scheme: "file",
			Path:   "/tmp/f1",
		}).Error; err != nil {
		t.Fatalf("failed to find record with url, got error: %v", err)
	}
	AssertEqual(t, uf1.String(), result.Storage.String())

	result = StructWithURL{}
	if err := DB.First(&result, "file_name = ? AND storage = ?", "FRemote2",
		datatypes.URL{
			Scheme:      "sftp",
			User:        url.UserPassword("user", "pwd"),
			Host:        "127.0.0.1",
			Path:        "/f2",
			RawPath:     "should not affects",
			RawQuery:    "query=1",
			Fragment:    "frag",
			RawFragment: "should not affects",
		}).Error; err != nil {
		t.Fatalf("failed to find record with url, got error: %v", err)
	}
	AssertEqual(t, u2.String(), uf2.String())
	AssertEqual(t, uf2.String(), result.Storage.String())
	AssertEqual(t, us, result.Storage.String())

	result = StructWithURL{}
	if err := DB.First(&result, "file_name = ? AND storage = ?", "FRemote2",
		datatypes.URL{
			Scheme:   "sftp",
			Opaque:   "//user:pwd@127.0.0.1/f2",
			RawQuery: "query=1",
			Fragment: "frag",
		}).Error; err != nil {
		t.Fatalf("failed to find record with url, got error: %v", err)
	}
	AssertEqual(t, us, result.Storage.String())

	result = StructWithURL{}
	if err := DB.First(&result, "file_name = ? AND storage = ?", "FRemote2",
		datatypes.URL{
			Scheme:   "sftp",
			User:     url.User("user"),
			Host:     "127.0.0.1",
			Path:     "/f2",
			RawQuery: "query=1",
			Fragment: "frag",
		}).Error; err == nil {
		t.Fatalf("record couldn't have been identical: %v vs %v", result.Storage, us)
	}
}
