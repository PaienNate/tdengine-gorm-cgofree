package tdengine_gorm

import (
	"strings"
	"testing"

	_ "github.com/taosdata/driver-go/v3/taosWS"
	"gorm.io/gorm"
)

func TestTaosWSRejectsTCPDSNBeforeConnecting(t *testing.T) {
	db, err := gorm.Open(&Dialect{
		DriverName: "taosWS",
		DSN:        "user:secret@tcp(example-host:6030)/gorm_test?loc=Local",
	}, &gorm.Config{})

	if err == nil {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		t.Fatal("expected taosWS with tcp DSN to fail before opening connection")
	}

	msg := err.Error()
	if !strings.Contains(msg, "taosWS requires ws/wss DSN") {
		t.Fatalf("expected actionable taosWS DSN error, got %q", msg)
	}
	if !strings.Contains(msg, "got tcp") {
		t.Fatalf("expected error to include actual network, got %q", msg)
	}
	if strings.Contains(msg, "secret") {
		t.Fatalf("expected error not to leak password, got %q", msg)
	}
}

func TestTaosWSAcceptsWebSocketDSN(t *testing.T) {
	db, err := gorm.Open(&Dialect{
		DriverName: "taosWS",
		DSN:        "user:secret@ws(example-host:6041)/gorm_test?loc=Local",
	}, &gorm.Config{})
	if err != nil {
		t.Fatalf("expected taosWS with ws DSN to pass local validation, got %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("expected sql DB handle, got %v", err)
	}
	_ = sqlDB.Close()
}

func TestDefaultDriverNameStillControlsImplicitDriverSelection(t *testing.T) {
	original := DefaultDriverName
	DefaultDriverName = "taosWS"
	defer func() {
		DefaultDriverName = original
	}()

	db, err := gorm.Open(Open("user:secret@tcp(example-host:6030)/gorm_test?loc=Local"), &gorm.Config{})
	if err == nil {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		t.Fatal("expected implicit driver selection to honor taosWS validation")
	}

	msg := err.Error()
	if !strings.Contains(msg, "taosWS requires ws/wss DSN") {
		t.Fatalf("expected taosWS validation error, got %q", msg)
	}
}
