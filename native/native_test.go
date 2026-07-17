//go:build cgo

package native

import (
	"testing"

	"gorm.io/gorm"
)

func TestOpenUsesTaosSqlDriver(t *testing.T) {
	db, err := gorm.Open(Open("user:secret@tcp(example-host:6030)/gorm_test?loc=Local"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected native.Open to initialize dialect, got %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("expected sql DB handle, got %v", err)
	}
	_ = sqlDB.Close()
}
