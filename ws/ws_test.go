package ws

import (
	"testing"

	"gorm.io/gorm"
)

func TestOpenUsesTaosWSDriver(t *testing.T) {
	db, err := gorm.Open(Open("user:secret@ws(example-host:6041)/gorm_test?loc=Local"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected ws.Open to initialize dialect, got %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("expected sql DB handle, got %v", err)
	}
	_ = sqlDB.Close()
}
